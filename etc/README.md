| 模块名称     | REST端口 | RPC端口 | Prometheus |
|----------|--------|-------|------------|
| system   | 8090   | 9091  | 4001       |
| monitor  | 8091   | 9092  | 4002       |
| resource | 8092   | 9093  | 4003       |

### 新建模块
```shell
goctl api new demo --style go_zero
```