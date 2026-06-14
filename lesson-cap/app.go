package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/kbinani/screenshot"
	"github.com/signintech/gopdf"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ROI struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type App struct {
	ctx           context.Context
	roi           *ROI
	isCapturing   bool
	tempDir       string
	capturedFiles []string
	mu            sync.Mutex
	stopChan      chan struct{}
}

func NewApp() *App {
	return &App{
		stopChan: make(chan struct{}),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dir, _ := os.MkdirTemp("", "lesson-cap-*")
	a.tempDir = dir
}

// CheckPermission: Detects if Screen Recording is enabled
func (a *App) CheckPermission() bool {
	// Try to capture a small patch at the center
	img, err := screenshot.CaptureRect(image.Rect(100, 100, 110, 110))
	if err != nil || img == nil { return false }
	
	// Check if we actually got pixels (blocked apps return transparent)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 { return true }
		}
	}
	return false
}

func (a *App) OpenSystemSettings() {
	exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_ScreenCapture").Run()
}

// SelectArea V36: Reverted to the SUCCESSFUL V22 logic (No delays, Step 2 scan)
func (a *App) SelectArea() (string, error) {
	// 1. Capture screen reference IMMEDIATELY (V22 style)
	bounds := screenshot.GetDisplayBounds(0)
	fullScreen, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "无法截图，请确认权限", fmt.Errorf("PERM")
	}

	// 2. Hide window and start selection
	runtime.WindowHide(a.ctx)
	defer runtime.WindowShow(a.ctx)

	tmpPath := filepath.Join(a.tempDir, "select.png")
	os.Remove(tmpPath)

	cmd := exec.Command("screencapture", "-i", "-x", "-o", tmpPath)
	if err := cmd.Run(); err != nil {
		return "已取消", nil
	}

	// 3. Load user crop
	f, _ := os.Open(tmpPath)
	cropped, _, _ := image.Decode(f)
	f.Close()
	if cropped == nil { return "截图识别失败", nil }

	// 4. Use the proven V22 findSubImage algorithm
	foundX, foundY := findSubImageV22(fullScreen, cropped)
	
	if foundX != -1 {
		roi := &ROI{
			X:      foundX + bounds.Min.X,
			Y:      foundY + bounds.Min.Y,
			Width:  cropped.Bounds().Dx(),
			Height: cropped.Bounds().Dy(),
		}
		a.mu.Lock()
		a.roi = roi
		a.mu.Unlock()
		return fmt.Sprintf("锁定成功: %dx%d", roi.Width, roi.Height), nil
	}

	return "锁定失败：请不要在视频播放时拉框", nil
}

func findSubImageV22(full, sub image.Image) (int, int) {
	fullB, subB := full.Bounds(), sub.Bounds()
	subW, subH := subB.Dx(), subB.Dy()
	fullW, fullH := fullB.Dx(), fullB.Dy()
	if subW > fullW || subH > fullH { return -1, -1 }

	// V22 exact logic: Check middle and top-left with Step 2
	midX, midY := subW/2, subH/2
	cMid := sub.At(midX, midY)
	cTL := sub.At(0, 0)

	for y := 0; y <= fullH-subH; y += 2 {
		for x := 0; x <= fullW-subW; x += 2 {
			if colorsMatch(full.At(x+midX, y+midY), cMid) {
				if colorsMatch(full.At(x, y), cTL) {
					return x, y
				}
			}
		}
	}
	return -1, -1
}

func colorsMatch(c1, c2 color.Color) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	const threshold = 3000 // High tolerance
	absDiff := func(a, b uint32) uint32 { if a > b { return a - b }; return b - a }
	return absDiff(r1, r2) < threshold && absDiff(g1, g2) < threshold && absDiff(b1, b2) < threshold
}

func (a *App) StartCapture() string {
	a.mu.Lock()
	if a.roi == nil { a.mu.Unlock(); return "请先选择区域" }
	if a.isCapturing { a.mu.Unlock(); return "运行中" }
	a.isCapturing = true
	a.mu.Unlock()
	go a.captureLoop()
	return "Capture started."
}

func (a *App) StopCapture() string {
	a.mu.Lock(); defer a.mu.Unlock()
	if !a.isCapturing { return "未开始" }
	a.isCapturing = false
	a.stopChan <- struct{}{}
	return fmt.Sprintf("停止。共截取 %d 张。", len(a.capturedFiles))
}

func (a *App) captureLoop() {
	var lastHash *goimagehash.ImageHash
	ticker := time.NewTicker(1200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-a.stopChan: return
		case <-ticker.C:
			a.mu.Lock(); roi := a.roi; a.mu.Unlock()
			if roi == nil { continue }
			img, err := screenshot.CaptureRect(image.Rect(roi.X, roi.Y, roi.X+roi.Width, roi.Y+roi.Height))
			if err != nil { continue }
			curHash, _ := goimagehash.DifferenceHash(img)
			if lastHash == nil { 
				a.saveImage(img)
				lastHash = curHash
				runtime.EventsEmit(a.ctx, "new_slide", len(a.capturedFiles))
				continue 
			}
			dist, _ := curHash.Distance(lastHash)
			if dist > 6 {
				a.saveImage(img)
				lastHash = curHash
				runtime.EventsEmit(a.ctx, "new_slide", len(a.capturedFiles))
			}
		}
	}
}

func (a *App) saveImage(img image.Image) {
	a.mu.Lock(); defer a.mu.Unlock()
	path := filepath.Join(a.tempDir, fmt.Sprintf("s_%d.png", len(a.capturedFiles)))
	f, _ := os.Create(path); png.Encode(f, img); f.Close()
	a.capturedFiles = append(a.capturedFiles, path)
}

func (a *App) ExportPDF(outputPath string) string {
	a.mu.Lock(); files := make([]string, len(a.capturedFiles)); copy(files, a.capturedFiles); a.mu.Unlock()
	if len(files) == 0 { return "无内容" }
	pdf := gopdf.GoPdf{}; pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for _, file := range files {
		f, _ := os.Open(file); img, _, _ := image.DecodeConfig(f); f.Close()
		rect := &gopdf.Rect{W: float64(img.Width), H: float64(img.Height)}
		pdf.AddPageWithOption(gopdf.PageOption{PageSize: rect})
		pdf.Image(file, 0, 0, rect)
	}
	if outputPath == "" { outputPath = filepath.Join(os.Getenv("HOME"), "Desktop", "Lesson_"+time.Now().Format("150405")+".pdf") }
	pdf.WritePdf(outputPath); return "PDF已保存至桌面"
}

func (a *App) GetCapturedCount() int { a.mu.Lock(); defer a.mu.Unlock(); return len(a.capturedFiles) }
func (a *App) OpenFolder() { exec.Command("open", a.tempDir).Start() }
func (a *App) Quit() { runtime.Quit(a.ctx) }

// Deprecated stubs
func (a *App) ConfirmROI(roi ROI) {}
func (a *App) ConfirmSnip(roi ROI) {}
func (a *App) EnterFrameMode() {}
func (a *App) ExitFrameMode() ROI { return ROI{} }
func (a *App) GetMirrorData() (map[string]interface{}, error) { return nil, nil }
func (a *App) GetSnipData() (map[string]interface{}, error) { return nil, nil }
