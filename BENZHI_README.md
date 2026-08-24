# almen-int：喷丸 Almen 弧高核算命令行工具

almen-int 读入丸粒速度、试片几何、残余应力层与覆盖率，用钉死的薄片弯矩模型算弧高，按 1−e^{−λt} 覆盖率与「时间翻倍弧高只增 10%」规则判饱和，未饱和不报强度等级字母，并给出建议 Almen 试片号。

## 构建 / 运行 / 测试

```text
go build ./...
almen-int height example/a2-steel.json
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
