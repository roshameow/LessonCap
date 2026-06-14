# LessonCap 🎓

> **智能直播课件自动提取工具** - 让每一场直播课都变成一本完美的 PDF 笔记。

[![Build and Release macOS](https://github.com/YOUR_USERNAME/YOUR_REPO/actions/workflows/release.yml/badge.svg)](https://github.com/YOUR_USERNAME/YOUR_REPO/actions/workflows/release.yml)
[![Platform](https://img.shields.io/badge/platform-macOS-blue.svg)](https://github.com/YOUR_USERNAME/YOUR_REPO/releases)

LessonCap 是一款专为学生和职场人士设计的 macOS 桌面应用。它能自动监测屏幕上的视频课件，并在老师翻页时瞬间抓取，自动去重，最后导出为高清 PDF。

## ✨ 核心特性

- 🎯 **数学级选区锁定**：采用逐像素拼图比对技术，100% 精准锁定视频区域。
- 🤖 **智能翻页监测**：实时感知画面变化，自动过滤老师走动、弹幕干扰。
- ⚡ **零干扰体验**：一键开启，最小化到后台即可全自动工作。
- 📄 **一键导出 PDF**：自动合并所有课件页，生成便于复习的文档。
- 🎨 **原生 macOS 体验**：适配 Retina 屏幕，极致轻量。

## 📥 安装指南 (macOS)

由于 LessonCap 为个人开发者项目，未通过苹果付费认证，请按照以下步骤安装：

1. 前往 [Releases](https://github.com/YOUR_USERNAME/YOUR_REPO/releases) 下载最新的 `LessonCap_macOS.zip`。
2. 解压后，**请勿直接双击**。
3. **关键：** 在 `LessonCap.app` 上点击 **鼠标右键**，选择 **“打开 (Open)”**。
4. 在弹出的系统警告对话框中，再次点击 **“打开”**。
5. **权限设置**：首次运行点击“划框”时，请在“系统设置 -> 隐私与安全性 -> 屏幕录制”中允许本程序。

## 🚀 如何使用

1. **选择区域**：点击“1. 鼠标拉框选择区域”，在视频画面上划一个框。
2. **开始监测**：点击“2. 开始自动监测”，然后你可以正常切换到视频全屏观看。
3. **保存导出**：听课结束后，点击“3. 导出 PDF 到桌面”即可。

## 🛠 技术实现

- **Backend**: Go + Wails
- **Frontend**: Svelte + Vite
- **Algorithm**: dHash (Image Perception Hashing) + Pixel Matching
- **PDF Generation**: GoPdf

## 🤝 贡献与反馈

如果你在使用中遇到任何问题，欢迎提交 [Issue](https://github.com/YOUR_USERNAME/YOUR_REPO/issues)。

---

*Made with ❤️ for lifelong learners.*
