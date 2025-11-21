module.exports = {
  apps: [
    {
      name: "gin-template",
      script: "/home/admin/gin-template/main",
      args: "-c config/main.yml",
      cwd: "/home/admin/gin-template",

      // 进程管理
      instances: 1,
      exec_mode: "fork",
      autorestart: true,
      watch: false,
      max_memory_restart: "1G",

      // 重启策略
      restart_delay: 5000, // 对应 RestartSec=5s
      min_uptime: "10s",
      max_restarts: 10,

      // 日志配置
      out_file: "/home/admin/gin-template/std.log",
      error_file: "/home/admin/gin-template/std.log",
      merge_logs: true,
      log_date_format: "YYYY-MM-DD HH:mm:ss Z",

      // 环境变量
      env: {
        NODE_ENV: "production",
      },

      // 用户配置 (需要 PM2 以 root 运行才能切换用户)
      // uid: 'admin',
      // gid: 'admin',
    },
  ],
};
