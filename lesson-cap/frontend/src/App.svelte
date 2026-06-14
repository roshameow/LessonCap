<script>
  import { onMount } from 'svelte';
  import { 
    StartCapture, 
    StopCapture, 
    SelectArea,
    ExportPDF, 
    GetCapturedCount,
    OpenFolder
  } from '../wailsjs/go/main/App';
  import { EventsOn, Quit } from '../wailsjs/runtime/runtime';

  let status = "准备就绪";
  let capturedCount = 0;
  let isCapturing = false;
  let roiText = "未选择区域";

  onMount(async () => {
    capturedCount = await GetCapturedCount();
    
    EventsOn("new_slide", (count) => {
      capturedCount = count;
      status = `检测到课件翻页，已自动捕获！(第 ${count} 张)`;
    });
  });

  async function handleSelect() {
    status = "正在启动选择工具...";
    try {
      const res = await SelectArea();
      roiText = res;
      if (res.includes("锁定成功")) {
          status = "区域选择成功！";
      } else {
          status = res;
      }
    } catch (e) {
      status = "选区失败，请重试";
    }
  }

  async function toggleCapture() {
    if (isCapturing) {
      status = await StopCapture();
      isCapturing = false;
    } else {
      let res = await StartCapture();
      if (res === "Capture started.") {
        isCapturing = true;
        status = "正在监测中... (第 1 张已存)";
      } else {
        status = res;
      }
    }
  }

  async function handleExport() {
    status = "正在导出 PDF 到桌面...";
    const res = await ExportPDF("");
    status = res;
  }
</script>

<main>
  <div class="container">
    <header>
      <div class="title-bar">
        <h1>LessonCap 🎓</h1>
        <button class="quit-btn" on:click={Quit}>✕</button>
      </div>
      <p class="subtitle">智能直播课件自动提取 · V36 最终版</p>
    </header>

    <div class="stats">
      <div class="stat-item">
        <span class="label">当前状态</span>
        <span class="value status-text">{status}</span>
      </div>
      <div class="stat-item">
        <span class="label">锁定区域</span>
        <span class="value roi-info">{roiText}</span>
      </div>
    </div>

    <div class="controls">
      <button class="btn secondary" on:click={handleSelect}>
        🎯 1. 鼠标拉框选择区域
      </button>

      <button class="btn primary" class:active={isCapturing} on:click={toggleCapture}>
        {isCapturing ? "⏹ 停止自动监测" : "▶️ 2. 开始自动监测"}
      </button>

      <button class="btn success" on:click={handleExport} disabled={isCapturing || capturedCount === 0}>
        📄 3. 导出 PDF 到桌面
      </button>
    </div>

    <div class="hint-box">
        💡 <b>重要：</b> 如果点击按钮无反应，请前往：<br/>
        [系统设置] -> [隐私与安全性] -> [屏幕录制]，<br/>
        将 <b>LessonCap</b> 的开关<b>关闭再重新开启</b>。
    </div>

    <div class="footer">
       <button class="text-link" on:click={OpenFolder}>查看临时截图</button>
       <div class="count-box">已截取: {capturedCount}</div>
    </div>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: 'PingFang SC', -apple-system, sans-serif;
    background-color: #1b2636;
    color: white;
  }

  .container {
    padding: 24px;
    height: 100vh;
    display: flex;
    flex-direction: column;
    box-sizing: border-box;
  }

  .title-bar { display: flex; justify-content: space-between; align-items: center; }
  h1 { margin: 0; font-size: 1.8rem; color: #64ffda; }
  .quit-btn { background: none; border: none; color: #8892b0; font-size: 1.5rem; cursor: pointer; }
  .subtitle { color: #8892b0; margin: 5px 0 25px; font-size: 0.9rem; }

  .stats { background: #233554; padding: 15px; border-radius: 8px; margin-bottom: 25px; }
  .stat-item { margin-bottom: 12px; }
  .label { font-size: 0.75rem; color: #8892b0; display: block; margin-bottom: 4px; }
  .value { font-size: 0.9rem; font-weight: 500; }
  .status-text { color: #64ffda; }
  .roi-info { font-family: monospace; color: #64ffda; }

  .controls { display: flex; flex-direction: column; gap: 12px; }
  .btn { padding: 12px; border: none; border-radius: 6px; font-size: 1rem; font-weight: 600; cursor: pointer; transition: 0.2s; }
  .primary { background: #64ffda; color: #0a192f; }
  .primary.active { background: #ff4d4d; color: white; }
  .secondary { background: #334155; color: white; }
  .success { background: #4CAF50; color: white; }
  .success:disabled { background: #1e293b; color: #7f8c8d; cursor: not-allowed; }

  .hint-box { margin-top: 15px; font-size: 0.7rem; color: #ffbc2d; text-align: center; opacity: 0.8; line-height: 1.5; background: rgba(0,0,0,0.2); padding: 10px; border-radius: 6px; }

  .footer { margin-top: auto; display: flex; justify-content: space-between; align-items: center; font-size: 0.8rem; }
  .text-link { color: #64ffda; text-decoration: underline; background: none; border: none; cursor: pointer; padding: 0; }
  .count-box { background: rgba(255,255,255,0.1); padding: 2px 10px; border-radius: 12px; color: #8892b0; }
</style>
