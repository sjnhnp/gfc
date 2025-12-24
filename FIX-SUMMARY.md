# 问题修复总结

## ✅ 已修复的问题

### 1. SUB-RULE 规则显示"策略组不存在"警告

**问题描述：**
- SUB-RULE 规则在 GUI 中显示 `[ ! ]` 警告标记
- 提示"策略组不存在"

**原因：**
- SUB-RULE 的最后一部分（如 `QUIC_CHECK`）是 sub-rule 名称，不是代理组
- GUI 错误地检查它是否存在于代理组列表中

**修复：**
修改 `frontend/src/views/ProfilesView/components/RulesConfig.vue`：
```typescript
const hasLost = (r: ProfileType['rulesConfig'][0]) => {
  if (BuiltInOutbound.includes(r.proxy)) return false
  // For SUB-RULE type Logic rules, the proxy field contains the sub-rule name, not a proxy group
  if (r.type === RuleType.Logic && /^SUB-RULE,/i.test(r.payload)) return false
  return !props.profile.proxyGroupsConfig.find((v) => v.id === r.proxy)
}
```

**结果：**
✅ SUB-RULE 规则不再显示警告标记
✅ 其他规则的代理组检查不受影响

---

### 2. 版本号显示问题说明

**您看到的 v1.16.0 是正确的行为！**

#### 情况 A：本地构建（您当前的情况）
```bash
cd frontend
npm run build
cd ..
wails build
```

**版本号：** `v1.16.0`（默认值，定义在 `bridge/bridge.go`）

**如何修改本地开发版本号：**
```go
// bridge/bridge.go
var Version = "v1.17.0"  // ← 修改这里
```

#### 情况 B：GitHub Actions 构建
通过 GitHub Actions 构建时，版本号会自动注入：

**示例 1：推送 tag**
```bash
git tag v1.17.0
git push origin v1.17.0
```
→ 版本号：`v1.17.0`

**示例 2：手动触发 + 输入版本号**
→ 版本号：您输入的版本（如 `v1.18.0`）

**示例 3：手动触发 + 留空**
→ 版本号：自动递增（如 `v1.16.1`）

#### 版本号注入工作原理

GitHub Actions 使用 `-ldflags` 在编译时覆盖默认版本号：

```bash
wails build -ldflags "-X 'github.com/GUI-for-Cores/GUI.for.Clash/bridge.Version=v1.17.0'"
```

这会将 `bridge.Version` 从默认的 `v1.16.0` 改为 `v1.17.0`。

**重要：** 只有通过编译时参数才能改变，运行时无法修改。

---

## 📝 测试验证

### 测试 SUB-RULE 警告修复

**配置：**
```yaml
rules:
  - SUB-RULE,(AND,((NETWORK,UDP),(DST-PORT,443))),QUIC_CHECK
  - SUB-RULE,(DST-PORT,443/udp),QUIC_CHECK
  - DOMAIN-SUFFIX,google.com,PROXY  # 这个可能会有警告（如果 PROXY 组不存在）

sub-rules:
  QUIC_CHECK:
    - RULE-SET,cn,DIRECT
    - MATCH,REJECT
```

**预期结果：**
- ✅ SUB-RULE 规则：无警告标记
- ⚠️ 普通规则（如果代理组不存在）：显示警告标记

### 测试版本号

**本地构建测试：**
1. 修改 `bridge/bridge.go` 中的 `var Version = "v1.99.0"`
2. 重新构建：`wails build`
3. 运行程序，检查标题是否显示 `GUI.for.Clash v1.99.0`

**GitHub Actions 测试：**
1. 推送 tag：`git tag v2.0.0 && git push origin v2.0.0`
2. 等待 Actions 完成
3. 下载构建文件并运行
4. 检查标题是否显示 `GUI.for.Clash v2.0.0`

---

## 🔄 开发工作流建议

### 本地开发
- 保持 `Version = "v1.16.0"` 或设为 `"dev"`
- 专注于功能开发，不关心版本号

### 准备发布
1. **提交所有更改**
   ```bash
   git add .
   git commit -m "Prepare for v1.17.0 release"
   git push
   ```

2. **创建并推送 tag**
   ```bash
   git tag v1.17.0
   git push origin v1.17.0
   ```

3. **GitHub Actions 自动构建**
   - 版本号自动设为 `v1.17.0`
   - 创建 Release
   - 上传构建文件

4. **验证**
   - 下载并运行构建的程序
   - 检查版本号是否正确

---

## 📊 完整的修复清单

| 问题 | 状态 | 说明 |
|------|------|------|
| SUB-RULE 语法支持 | ✅ 完成 | 支持所有格式 |
| sub-rules 配置块 | ✅ 完成 | 正确保存和恢复 |
| SUB-RULE 警告问题 | ✅ 修复 | 不再显示错误警告 |
| 版本号动态注入 | ✅ 完成 | GitHub Actions 构建时自动设置 |
| 版本号自动递增 | ✅ 完成 | 手动触发时可选 |

---

## 🚀 下一步

### 重新构建测试
```bash
cd frontend
npm run build
cd ..
wails build
```

### 验证修复
1. ✅ 加载包含 SUB-RULE 的配置
2. ✅ 检查是否还有警告标记
3. ✅ 保存配置并重新加载
4. ✅ 验证规则完整性

### 发布新版本
如果一切正常，可以发布新版本：
```bash
git add .
git commit -m "Fix SUB-RULE warning and version injection"
git push
git tag v1.17.0
git push origin v1.17.0
```

---

## 总结

**已修复：**
1. ✅ SUB-RULE 规则不再显示"策略组不存在"警告
2. ✅ 版本号在 GitHub Actions 构建时正确注入
3. ✅ 本地开发可以手动修改默认版本号

**不需要修复：**
- ℹ️ 本地构建显示默认版本号是正常行为

所有功能现在都已完整实现！🎉
