Typecho 应用商店 Server
================

![Typecho应用商店][1]

## 实现原理

目前 Typecho 的插件集中在下面的两个仓库

- https://github.com/typecho/plugins
- https://github.com/chekun/typecho-fans-plugins

我使用go写了个服务，每一个小时去更新这2个仓库的数据，并解析插件文件的Meta信息入库打包。

## 使用方法

很简单，下载配套的的[AppStore插件][2]即可使用：


## 如何开发插件

- 遵守Typecho开发规范
    - 插件入口文件是Plugin.php，注意: **大写首字母**
    - 设置插件信息，目前采用的标签:
        * @package 
        * @author
        * @link
        * @version
        * @dependence
- 上传到Github
    - 使用git subtree将插件更新到typecho/plugins或者typecho-fans-plugins上
    - 注意插件版本控制，如果你改了代码，但版本号不变，那么视为没有更新。

这样你的插件就可以出现在了AppStore中了，当然这不是时时显示的呢。

> 如果以上两个项目的维护者可以让我设置个钩子最好了，那么我就不需要设置cronjob也会实时同步插件了。

## 欢迎一起完善该项目！

  [1]: http://chekun-blog.qiniudn.com/typecho-app-store-splash.png
  [2]: https://github.com/typecho-app-store/AppStore
