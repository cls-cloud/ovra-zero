| 模块名称     | REST端口 | RPC端口 | Prometheus |
|----------|--------|-------|------------|
| auth     | 8091   | -     | 4001       |
| system   | 8092   | 9092  | 4002       |
| demo     | 8099   | 9099  | 4009       |

### 新建模块
```shell
goctl api new demo --style go_zero
```