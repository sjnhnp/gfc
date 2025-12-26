<div align="center">
  <img src="build/appicon.png" alt="GUI.for.Clash" width="200">
  <h1>GUI.for.Clash (Enhanced Fork)</h1>
  <p>基于 <a href="https://github.com/GUI-for-Cores/GUI.for.Clash">GUI-for-Cores/GUI.for.Clash</a> 的增强版 Fork</p>
  <p>A GUI program for mihomo (Clash.Meta) kernel, developed with Vue3 + Wails.</p>
</div>

## ✨ Fork 增强功能

### 🔧 Bug 修复

| 修复项 | 说明 |
|--------|------|
| **SUB-RULE 语法支持** | 正确支持 Clash Mihomo SUB-RULE 语法，修复"策略组不存在"警告 |
| **自定义 GeoIP 配置** | 修复 `geodata-mode: false` 时 GeoIP 配置不生效的问题 |
| **FakeIP Filter 规则集** | 修复 `fake-ip-filter` 中使用 `rule-set` 时内核报错的问题 |
| **Mixin 规则集** | 在 mixin 中定义的规则集现在会正确复制到生成的配置中 |
| **DNS 规则集显示** | 仅用于 DNS 的规则集不再错误显示在 GUI 规则列表中 |
| **本地订阅恢复** | 本地订阅导入后配置能正确恢复 |

### 🔌 预装增强版 Gists 同步插件

- **一键备份/恢复**：支持 GitHub Gists 备份所有配置
- **本地文件支持**：
  - 备份本地订阅源文件（`data/local/*.txt`）
  - 备份本地规则集源文件
  - 备份插件第三方依赖（`data/third/`）
- **开箱即用**：插件和 `crypto-js.js` 依赖已预打包，无需手动安装

### 🛡️ Windows 优化

- **默认管理员权限**：Windows 版本默认以管理员身份运行
- **TUN 模式开箱即用**：无需手动设置权限

### ⚙️ 配置优化

- **系统代理默认关闭**：避免首次启动自动设置系统代理
- **增加自定义GitHub加速地址**：解决第一次运行无法下载内核

### � 构建优化

- **GitHub Actions 自动构建**：支持 Windows 和 macOS 自动构建发布
- **版本号自动递增**：每次发布自动递增 patch 版本号
- **启动时间优化**：插件预打包，减少首次启动的网络请求


## 📖 文档

- [使用指南](https://gui-for-cores.github.io/guide/gfc/how-to-use)


## 🙏 致谢

- [GUI-for-Cores/GUI.for.Clash](https://github.com/GUI-for-Cores/GUI.for.Clash) - 原项目
- [MetaCubeX/mihomo](https://github.com/MetaCubeX/mihomo) - Clash.Meta 内核

## 📄 License

MIT License - 详见 [LICENSE](LICENSE)
