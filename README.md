# Meta-Front-G（Antd Pro 代码生成器）

## 作用

根据`typings.d.ts`生成`Column`、`Page`、`routes.ts`文件

## 为什么会有这个项目

Antd Pro 项目可以通过 openAPI 文档生成 Services 文件，但是还是需要自己一个个去写页面，有点麻烦，虽然 pro 组件提取了通用的组件，但是组合起来用还是会增加不少工作量（写代码、需要翻文档等）。

于是诞生了这个项目，通过手动提取各个相关组件，完成通用页面后，根据 `typings.d.ts`文件生成列数据、页面数据、路由文件；后续只需调整字段，调整新增/编辑表单，处理字段数据转换，调整路由即可完成页面开发。

## Installation

`go install github.com/one-meta/meta-front-g@latest`

或 clone 项目后，`go install`

## 用法（步骤）【项目使用建议参考 [meta wiki](https://github.com/one-meta/meta/wiki/)】

1. 复制`typings.d.ts`到生成器目录，或者在运行时指定 `typings.d.ts`文件的绝对路径
2. 根据实际情况修改页面模板 BasePage 中的页面（也可以生成后，再修改生成的页面）
3. 运行生成器，会在当前目录生成 Column：列，Pages：页面，routes.ts：路由文件
4. 复制需要的列到前端框架的`src/columns`下，根据实际情况修改列属性，或增删列，增加导入
5. 复制需要的页面文件夹到前端框架的`src/pages`下，根据实际情况调整新增/编辑的表单，或调整其他功能，增加导入
6. 复制`routes.ts`中所需要的路由到前端框架的 `config/routes.ts`
7. 通过以上步骤，就能快速生成基于通用组件的增删改查页面

## 更多信息

### config.toml

生成器配置文件，用于配置字段生成 antd pro column 的模板，包括：排除的实体名、排除的字段、扩展字段、解析方式及模板等

### typings.d.ts

命名空间下声明了各个实体的类型

在 antd pro 项目 npm run openapi 之后，会生成 services 文件和 typings.d.ts

[详情](https://pro.ant.design/zh-CN/docs/openapi)

大概意思：后端框架生成标准的 openAPI 文档（swagger api）之后，可以通过 npm run openapi 生成对应的 Service 文件，加快开发。
