# QQ空间相册下载
基于golang开发

## 特性
- 多线程下载
- 下载原始质量
- 下载相册中的视频
- 下载的文件修改时间为照片拍摄时间
- 文件名为照片在QQ相册中的ID

## 配置
`.env` 文件中的配置含义与获取方法
1. 打开浏览器调试工具，在Network选项卡中过滤 `cgi_list_photo`
2. 点进一个相册，这时候下面就会看到请求
3. 查看请求中的Cookie一栏，复制到 `COOKIESTR`
4. 查看url请求参数中的g_tk复制到 `G_TK`， 这应该是代表本次登录的ID
5. `QQ` 就是你当前要采集的QQ
6. `TOPICID` 是相册ID，也是在请求中获取
7. `DOWNLOADDIR` 就是你要下载到的路径，例如 `D:\Download\我的相册`

## TODO
- 可以选择相册
- 使用账号密码登录
