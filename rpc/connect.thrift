namespace go connect

// 定义服务
service ConnectionService {
    ConnectionResponse GetConnectionDetail()
}

// 定义返回的响应结构
struct ConnectionResponse {
    1: i32 chat;      // Chat 连接数量
    2: i32 example;   // Example 连接数量
    3: i32 num;       // 房间数量
}