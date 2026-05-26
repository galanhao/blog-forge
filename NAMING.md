# 命名规范 — blog-forge

本文件是所有代码的命名权威参考。写代码前先看这里。

## 核心领域词汇（全局统一，不可替换）

| 概念 | 统一命名 | ❌ 禁止使用 |
|------|---------|------------|
| 博客文章 | Post | article, entry, blog, document |
| 独立页面 | Page | standalone, special |
| 站点 | Site | website, blog, weblog |
| 主题 | Theme | skin, template, layout |
| 分类 | Category | section, group, folder |
| 标签 | Tag | label, keyword |
| 归档 | Archive | history, timeline |
| 摘要 | Excerpt | summary, abstract, description |
| 封面图 | CoverImage | thumbnail, hero, banner |
| 永久链接 | Permalink | url, link, path |
| Slug | Slug | url_slug, name, key |
| Frontmatter | Frontmatter | metadata, header, yaml_header |
| 配置 | Config | settings, options, prefs |
| 导航 | Nav | navigation, menu |
| 上下文 | Context | ctx (仅在函数参数中缩写), data, vars |
| 布局 | Layout | template_type, page_type |
| 资源 | Asset | resource, file, static_file |
| 静态文件 | Static | public, dist_file |
| 渲染 | Render | convert, transform, process |
| 构建 | Build | compile, generate, produce |
| 分页 | Pagination | pager, paging, pages |

## 包命名

- 全小写，单个词，无下划线
- ✅ `content`, `render`, `permalink`, `theme`, `site`
- ❌ `blogCore`, `post_service`, `templateEngine`

## 文件命名

- 全小写，下划线分隔
- 一个文件一个职责，文件名反映内容
- ✅ `frontmatter.go`, `permalink.go`, `pagination.go`
- ❌ `utils.go`, `helpers.go`, `misc.go`

## 函数命名

- Go 惯例：导出用 PascalCase，内部用 camelCase
- 函数名要做什么说什么，不要模糊
- ✅ `LoadPosts(dir string) ([]*Post, error)`
- ✅ `CalculatePermalink(format string, post *Post) string`
- ❌ `ProcessData()`
- ❌ `HandleStuff()`

## 结构体命名

- 领域概念直接用名词：`Post`, `Site`, `Theme`, `Category`
- 上下文结构体用 `XxxCtx` 后缀：`IndexCtx`, `PostCtx`, `PaginationCtx`
- 配置结构体用 `XxxConfig` 后缀：`SiteConfig`, `ThemeConfig`, `DeployConfig`
- 选项结构体用 `XxxOptions` 后缀

## 变量命名

- 短作用域用短名：`i`, `p`, `r`
- 长作用域用描述性名：`postData`, `themeConfig`
- 布尔变量用 `is/has/can/should` 前缀：`isPublished`, `hasCover`
- 一致性：同一概念在所有地方用同一个词

## 常量命名

- 导出常量 PascalCase：`LayoutPost`, `FormatDateFull`
- 内部常量 camelCase：`defaultPerPage`

## 错误信息

- 小写开头，不加标点
- ✅ `fmt.Errorf("post not found: %s", slug)`
- ❌ `fmt.Errorf("Post Not Found: %s", slug)`
- ❌ `fmt.Errorf("post not found: %s.", slug)`
