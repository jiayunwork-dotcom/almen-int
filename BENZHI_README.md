# almen-int: Go HTTP 后端核算喷丸 Almen 试片弧高与饱和强度

用户给出丸粒速度、试片几何、弹性模量、残余压应力层和覆盖率，核算薄片弯矩弧高、指数覆盖率增益、10% 时间加倍饱和判定，以及 N/A/C 试片等效强度。必须同时成立：速度加倍则动能与弧高按平方律上升；覆盖率从低到饱和弧高趋近平台且再加一倍时间弧高增幅不超过 10%；同一丸流加厚试片因 I∝t³ 使弧高下降。未饱和不得报强度等级字母。双面等增益时净弯矩为零；三角应力剖面在同一峰值下弯矩低于均匀层。

## 构建 / 运行 / 测试

```text
go build ./...
go run . serve
go run . height example/a2-steel.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
