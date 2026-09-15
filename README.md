# LessonCap

[下载安装包](https://github.com/roshameow/LessonCap/releases/latest) · [合成课件与 PDF 示例](docs/examples/README.md) · [问题反馈](https://github.com/roshameow/LessonCap/issues)

[![Build and Release](https://github.com/roshameow/LessonCap/actions/workflows/release.yml/badge.svg)](https://github.com/roshameow/LessonCap/actions/workflows/release.yml)

选择视频课件区域，监测页面变化，保存截图并导出 PDF。适合以静态幻灯片为主的直播课和视频课程。

## 三步使用

1. **选择区域**：框选课件内容，尽量避开老师画面、弹幕和播放器控件。
2. **开始监测**：保持课件区域可见，翻页后检查截图计数是否增加。
3. **停止并导出**：停止监测，再导出 PDF，检查页序、漏页和重复页。

首次试用可打开仓库中的 [slides.html](docs/examples/slides.html)，按说明逐页切换。也可查看 [PDF 格式示例](docs/examples/sample-slides.pdf)。示例使用合成内容，由脚本生成，**不是实际录屏导出结果或准确率证明**。

## 安装

| 系统 | 下载文件 | 安装说明 |
| --- | --- | --- |
| macOS | `LessonCap_macOS.zip` | 解压后将应用放到 Applications，首次运行需允许屏幕录制 |
| Windows | `LessonCap_Windows.zip` | 解压并运行 `LessonCap.exe` |

安装包见 [Releases](https://github.com/roshameow/LessonCap/releases/latest)。macOS 版本未公证，Windows 版本可能出现 SmartScreen 提示；系统具体提示以当前版本为准。macOS 可在“系统设置 → 隐私与安全性 → 屏幕录制”检查授权。

## 适用范围与限制

- 检测依据是图像变化，不是课件内容识别。动画、弹幕、光标和老师画面可能造成额外截图；相似页面也可能漏检。
- 当前实现约每 **1.2 秒**采样一次；快速翻页不保证全部捕获。
- 应用截取所选屏幕区域；被其他窗口遮挡、切换布局或移动播放器会影响结果。
- PDF 由截图组成，不包含 OCR 可搜索文本、语音转写或自动知识总结。
- 多屏、缩放和系统权限可能影响选区。先用合成课件检查本机效果，再用于完整课程。

## 从源码开发

技术栈：Go + Wails v2、Svelte + Vite、dHash 变化检测、GoPdf 导出。

安装 Go（版本见 `lesson-cap/go.mod`）、Node.js 22+ 及 [Wails 平台依赖](https://wails.io/docs/gettingstarted/installation/)，然后：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
cd lesson-cap
npm ci --prefix frontend
wails dev
```

仅检查前端：`npm run build --prefix lesson-cap/frontend`。这不会验证原生截图、权限或 PDF 导出。

## 更新与反馈

[CHANGELOG](CHANGELOG.md) 记录更新与发布说明。提交问题时请附系统版本、LessonCap 版本、显示器/缩放设置、复现步骤及合成示例。真实课件截图请先确认可以分享。
