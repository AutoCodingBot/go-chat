### 消息全部混在同张表
    初级:聊天表跟群聊表分开   ✔
    高级:分开,并用mongodb存

### LOG
    初级:log按daily记录
    高级:记录恶意用户;过于频繁的request,访问奇怪的路由
    
### 表UserMessage,GroupMessage 设计结构问题
    问题:拉取好友聊天信息,走的sql是senderId = X && receiverId = Y OR x,y互换
    Idea:新增一个字段 `conversation_id`,这个字段可以描述好友关系
        `hash` 比如 hash(username1 +username2) :碰撞结果未知,索引低效
        `int` 比如 100 +200 (+是concat,不是plus)

### 加密
    明文存储        ✔ 通过bcrypt加密

### jwt身份验证
    怎样的对称密钥才是性能和安全相平衡的
    示例代码:https://gitee.com/tang_q/gin-vue3-admin-api

### 图片缓存
    报文头已告知浏览器进行缓存(Cache-Controler  public,max-age=86400),但刷新页面后,前端图片依然从服务器获取
