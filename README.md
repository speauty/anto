## ANTO-不专业字幕翻译的Windows桌面应用

> 作为一名临时搬运工，搞个字幕翻译工具，一点也不过分~是吧
>



### 写在前面

我是21年底才接触到油管搬运。不知出于什么原因，看了UP主[@神经元猫](https://space.bilibili.com/364152971/?spm_id_from=333.999.0.0)搬运翻译的[Cherno C++ 中文](https://space.bilibili.com/364152971/channel/collectiondetail?sid=13909)，感觉[TheCherno](https://github.com/TheCherno)挺有趣的，就萌生了搬运他的OpenGL教程，然后就是游戏引擎等。刚开始，特别费人，使用剪映的智能字幕进行音频识别并导出SRT字幕。

由于分段很乱，需要人工调整，然后挨个翻译。遇到陌生单词，还要去查一下。但是吧，后来觉得效率低得令人发指。既然有智能字幕，那么就有智能翻译，然后，就了解到了各种机器翻译。看着每月的免费额度，完全可以操作一下。那么，就开始着手这款应用的研发。并不复杂，初版几天就搞定了，主要还是对框架不太熟悉，而且有些场景不太了解，大多都是后面慢慢完善。

用，是能用了，但是，应该还有很大提升的空间。比如，有朋友提到的AI翻译，以及通知和托盘等。

加油~

### 任务列表

- [x] 集成lxn/walk桌面GUI库（才把fyne清理了，那个编译运行太慢了，而且Win下，依赖环境有点。。。），这个要抽一波，不然用起来，小难受；
- [x] 实现静态界面（还是不大会多窗口模式，估计会照旧采用Widget的显示来模拟多页面，也可能还是单页）；
  - [x] 常用配置本地缓存（文件），就我最近的使用体验来看，不想再输密钥了；
  - [x] 增加全量和增量翻译模式区分，防止重复翻译，浪费资源；
  - [x] 增加批量翻译模式（挂机，一直跑），那个列表有点小烦，后面没有单独放置操作的列；
- [x] 实现srt文件解析和压缩（这块好搞，我感觉都写了好几次了）；
- [x] 接入上游翻译服务（如果有什么好的免费的服务，也可以给我说）；

### 应用预览

![应用启动](./assets/images/startup.jpg)

![设置界面](./assets/images/settings.jpg)

