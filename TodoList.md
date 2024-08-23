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

        前端载荷是:
        Uuid: 28353ed6-5966-4804-9c52-9b00abd4401e
        FriendUsername: eric
        MessageType: 1
        考虑 userName +friendName

        暂定方案
        names := []string{"eric", "sam"} // 创建一个包含两个字符串的切片

        // 使用sort.Strings函数对字符串切片进行排序
        sort.Strings(names)
        fmt.Println(names) // ["eric", "sam"]
        !!! username可能会被修改,这是不正确的行为

### 加密
    明文存储        ✔ 通过bcrypt加密

### jwt身份验证
    1.怎样的对称密钥才是性能和安全相平衡的

    2.jwt paylod存储用户ID风险评估
        安全风险
        信息泄露风险：payload中的信息（如用户ID）是明文存储的，任何人都可以查看到。这意味着，一旦JWT在网络中传输或存储不当（即使没有被篡改），用户ID就可能暴露给第三方。
        权限控制风险：如果系统仅仅依赖JWT中的用户ID来授权访问资源，恶意用户可能会尝试修改JWT的payload，尽管他们不能改变签名使JWT有效，但如果系统存在漏洞，比如未验证某些操作的权限，就可能利用已知的用户ID尝试访问不应访问的资源。
        令牌劫持：如果攻击者通过某种方式（如XSS、CSRF攻击）获取了用户的JWT，他们就能冒充该用户进行操作，特别是当JWT的有效期较长时，这种风险更加突出。
    业务风险
        隐私问题：用户可能对他们的ID被轻易地暴露感到担忧，尤其是在重视隐私的应用场景中。这可能影响用户信任度，导致合规性问题或用户流失。
        账户接管风险：结合其他漏洞，如弱口令或社会工程学攻击，暴露的用户ID可能帮助攻击者更容易地接管用户账户。
        审计和追踪难度：如果仅依赖JWT中的用户ID进行操作记录，一旦发生安全事件，追溯和分析攻击路径可能会较为困难，尤其是当攻击者尝试模仿合法用户行为时。

### 图片缓存
    报文头已告知浏览器进行缓存(Cache-Controler  public,max-age=86400),但刷新页面后,前端图片依然从服务器获取

### 在线用户信息采集
    新建常驻变量global/state/UserState (slice),

### 原生sql输出
    Debug()