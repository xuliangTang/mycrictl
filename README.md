## 配置 cni 插件
下载 [plugins](https://github.com/containernetworking/plugins) (v1.0.1) 压缩包

```bash
mkdir -p /opt/cni/bin/
mkdir -p /etc/cni/net.d
tar -zxvf cni-plugins-linux-amd64-v1.0.1.tgz  -C /opt/cni/bin/
```
从 [cni](https://github.com/containernetworking/cni) 文档抄一段配置
```bash
cat >/etc/cni/net.d/10-mynet.conf <<EOF
{
	"cniVersion": "0.2.0",
	"name": "mynet",
	"type": "bridge",
	"bridge": "cni0",
	"isGateway": true,
	"ipMasq": true,
	"ipam": {
		"type": "host-local",
		"subnet": "10.22.0.0/16",
		"routes": [
			{ "dst": "0.0.0.0/0" }
		]
	}
}
EOF
```
修改一个 containerd 配置
```bash
vi /etc/containerd/config.toml

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = false
```
重启 containerd
```bash
systemctl restart containerd
```
测试运行
```bash
cd cri
crictl run nginx.yaml sandbox.yaml

crictl pods
crictl ps
```