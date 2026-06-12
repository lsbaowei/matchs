# matchs

`matchs` 是一个用于文本关键词查找与替换的 Go 包。它面向“给定一批规则，判断输入文本命中了哪些规则，并按需返回替换后文本”的场景，支持普通关键词、组合规则和正则规则三类匹配方式。

## 安装

```bash
go get github.com/lsbaowei/matchs
```

## 主要功能

- **普通关键词匹配**：使用 `DFAMatcher` 构建关键词树，支持命中词返回和 rune 级替换。
- **组合规则匹配**：使用 `AssembleMatcher` 支持 `A|B|C` 顺序规则和 `A#B` 排除规则。
- **正则规则匹配**：使用 `RegexpMatcher` 支持以 `reg@` 开头的正则规则。
- **统一入口**：使用 `MatchService` 自动按规则类型分发，并按 `DFA -> ASSEMBLE -> REGEXP` 的稳定顺序执行。

## 快速开始

```go
package main

import (
	"fmt"

	"github.com/lsbaowei/matchs"
)

func main() {
	rules := []string{
		"敏感词",        // 普通关键词
		"A|B|C",       // A、B、C 依次出现
		"A#B",         // 出现 A 且不出现 B
		`reg@1\d{10}`, // 正则规则
	}

	service := matchs.NewMatchService()
	service.Build(rules)

	words, replaced := service.Match("敏感词 AxxBxxC 手机号13800138000", false, '*')

	fmt.Println(words)    // [敏感词 A|B|C reg@1\d{10}]
	fmt.Println(replaced) // *** AxxBxxC 手机号13800138000
}
```

## 规则格式

### 普通关键词

不包含 `|`、`#`，且不以 `reg@` 开头的规则会作为普通关键词处理。

```go
[]string{"敏感词", "违禁词"}
```

普通关键词支持替换文本返回：

```go
words, replaced := service.Match("这里有敏感词", false, '*')
// words:    []string{"敏感词"}
// replaced: "这里有***"
```

### 顺序组合规则：`|`

`A|B|C` 表示 `A`、`B`、`C` 必须按顺序出现在文本中。

```go
rules := []string{"A|B|C"}

// 命中
service.Match("AxxBxxC", false, '*')

// 不命中
service.Match("BxxAxxC", false, '*')
```

组合规则只返回命中的规则名，不执行脱敏替换。

### 排除组合规则：`#`

`A#B` 表示文本必须包含 `A`，且不能包含 `B`。

```go
rules := []string{"A#B"}

// 命中：包含 A，且不包含 B
service.Match("xxAxx", false, '*')

// 不命中：同时包含 A 和 B
service.Match("xxAxxBxx", false, '*')
```

### 正则规则：`reg@`

以 `reg@` 开头的规则会作为正则规则处理。返回的命中项是带 `reg@` 前缀的规则名。

```go
rules := []string{`reg@1\d{10}`}

words, replaced := service.Match("手机号13800138000", false, '*')
// words:    []string{`reg@1\d{10}`}
// replaced: "手机号13800138000"
```

正则规则只返回命中的规则名，不执行脱敏替换。非法正则会在 `Build` 时被跳过。

## API

### `MatchService`

```go
service := matchs.NewMatchService()
service.Build(rules)
words, replaced := service.Match(text, onlyOne, repl)
```

- `Build(rules []string)`：按规则格式自动分发到对应 matcher。重复调用会在现有 matcher 状态上追加规则；如果需要全量重建，请创建新的 `MatchService`。
- `Match(text string, onlyOne bool, repl rune)`：返回命中规则列表和替换后文本。
- `onlyOne=true`：按 `DFA -> ASSEMBLE -> REGEXP` 顺序返回首次命中结果。
- `repl`：普通关键词命中后用于替换的 rune，例如 `'*'`。

### 单独使用 matcher

也可以直接使用某一类 matcher：

```go
matcher := matchs.NewDFAMatcher()
matcher.Build([]string{"敏感词"})
words, replaced := matcher.Match("敏感词", false, '*')
```

可用 matcher：

- `NewDFAMatcher()`：普通关键词匹配与替换
- `NewAssembleMatcher()`：组合规则匹配
- `NewRegexpMatcher()`：正则规则匹配

历史拼写兼容入口仍保留，但已标记为 deprecated：

- `NewDFAMather()`
- `NewAssembleMather()`
- `AssembleMather`
- `RegexpMather`

新代码建议使用 `Matcher` 拼写正确的名称。

## 测试

```bash
go test ./...
```

## 注意事项

- 只有普通关键词 matcher 会修改替换文本；组合规则和正则规则会返回原文。
- 正则规则依赖 `github.com/dlclark/regexp2`。
- 规则分类由字符串格式决定：`reg@` 优先，其次 `|` / `#`，最后是普通关键词。
- 空普通关键词会被忽略，避免产生任意文本误命中。
