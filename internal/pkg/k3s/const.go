package k3s

const master = `[Unit]
Description=Lightweight Kubernetes
Documentation=https://k3s.io
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=notify
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/systemd/system/k3s-token.service.env
EnvironmentFile=-/etc/systemd/system/k3s-custom.service.env
KillMode=process
Delegate=yes
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/k3s server \
    --tls-san {{ .TLSSAN }} \
    {{- if .DataDir }}
    --data-dir {{ .DataDir }} \
    {{- end }}
		{{- if .ClusterCidr }}
    --cluster-cidr {{ .ClusterCidr }} \
    {{- end }}
		{{- if .ServiceCidr }}
    --service-cidr {{ .ServiceCidr }} \
    {{- end }}
    --service-node-port-range 20000-52767 \
    --token {{ .Token }} \
    {{- if not .Master0 }}
    --server {{ .Server }} \
    {{- end }}
    {{- if .Docker }}
    --docker \
    {{- end }}
    --cluster-init \
    --etcd-expose-metrics \
    --disable-network-policy \
    --disable-helm-controller \
    --disable servicelb,traefik \
    --kube-proxy-arg "proxy-mode=ipvs" "masquerade-all=true" \
    --kube-proxy-arg "metrics-bind-address=0.0.0.0"
`

const  agnet = `
[Unit]
Description=Lightweight Kubernetes
Documentation=https://k3s.io
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=notify
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/systemd/system/k3s-token.service.env
EnvironmentFile=-/etc/systemd/system/k3s-custom.service.env
KillMode=process
Delegate=yes
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/k3s agent \
    {{- if .DataDir }}
    --data-dir {{ .DataDir }} \
    {{- end }}
    --token {{ .Token }} \
    --server {{ .Server }} \
    {{- if .Docker }}
    --docker \
    {{- end }}
    --kube-proxy-arg "proxy-mode=ipvs" "masquerade-all=true" \
    --kube-proxy-arg "metrics-bind-address=0.0.0.0"
`
