好，现在进入 阶段三：接口与并发 🔥
  
  这个阶段难度明显提升了，我们要把 Java 的接口 + 线程池并发模型翻译成 Go 的隐式接口 + goroutine + channel 模型。先看这段 Java 代码：

  // ========== 接口定义 ==========
  public interface NotificationSender {
      void send(String message) throws SendException;
  }

  // ========== 自定义异常 ==========
  public class SendException extends Exception {
      public SendException(String message) {
          super(message);
      }
  }

  // ========== 邮件发送器 ==========
  public class EmailSender implements NotificationSender {
      private final String smtpServer;

      public EmailSender(String smtpServer) {
          this.smtpServer = smtpServer;
      }

      @Override
      public void send(String message) throws SendException {
          if (message == null || message.isEmpty()) {
              throw new SendException("邮件内容不能为空");
          }
          // 模拟发送延迟
          try {
              Thread.sleep(100);
          } catch (InterruptedException e) {
              Thread.currentThread().interrupt();
          }
          System.out.println("[" + smtpServer + "] 邮件已发送: " + message);
      }
  }

  // ========== 短信发送器 ==========
  public class SmsSender implements NotificationSender {
      private final String apiUrl;

      public SmsSender(String apiUrl) {
          this.apiUrl = apiUrl;
      }

      @Override
      public void send(String message) throws SendException {
          if (message == null || message.isEmpty()) {
              throw new SendException("短信内容不能为空");
          }
          // 模拟发送失败（故意让短信经常失败，用来测试错误收集）
          if (Math.random() < 0.3) {
              throw new SendException("短信网关超时: " + apiUrl);
          }
          System.out.println("[" + apiUrl + "] 短信已发送: " + message);
      }
  }

  // ========== 批量通知服务 ==========
  public class BatchNotifier {
      private final List<NotificationSender> senders;
      private final ExecutorService executor;

      public BatchNotifier(List<NotificationSender> senders) {
          this.senders = senders;
          this.executor = Executors.newFixedThreadPool(senders.size());
      }

      /**
       * 并发发送通知给所有发送器，收集所有错误
       * 返回发送失败的错误列表（空列表表示全部成功）
       */
      public List<String> sendAll(String message) {
          List<Future<Void>> futures = new ArrayList<>();
          List<String> errors = Collections.synchronizedList(new ArrayList<>());

          for (NotificationSender sender : senders) {
              futures.add(executor.submit(() -> {
                  try {
                      sender.send(message);
                  } catch (SendException e) {
                      errors.add(e.getMessage());
                  }
                  return null;
              }));
          }

          // 等待所有任务完成
          for (Future<Void> future : futures) {
              try {
                  future.get();
              } catch (InterruptedException | ExecutionException e) {
                  errors.add("任务执行异常: " + e.getMessage());
              }
          }

          return errors;
      }

      public void shutdown() {
          executor.shutdown();
      }
  }

  // ========== 主程序 ==========
  public class Main {
      public static void main(String[] args) {
          List<NotificationSender> senders = List.of(
              new EmailSender("smtp.example.com"),
              new EmailSender("smtp.backup.com"),
              new SmsSender("https://sms.api.com/send"),
              new SmsSender("https://sms.backup.com/send")
          );

          BatchNotifier notifier = new BatchNotifier(senders);

          // 测试正常发送
          System.out.println("=== 发送正常通知 ===");
          List<String> errors = notifier.sendAll("系统维护通知：今晚 22:00 进行例行维护");
      public static void main(String[] args) {
          List<NotificationSender> senders = List.of(
          if (errors.isEmpty()) {
              System.out.println("全部发送成功！");
          } else {
              System.out.println("部分发送失败：");
              errors.forEach(e -> System.out.println("  - " + e));
          }

          // 测试空消息
          System.out.println("\n=== 发送空消息 ===");
          List<String> emptyErrors = notifier.sendAll("");
          System.out.println("预期失败数: " + emptyErrors.size());

          notifier.shutdown();
      }
  }

  你的任务：请把上面这段 Java 代码改写成地道的 Go 代码 🎯

  关键考察点：
  1. NotificationSender 接口 → Go 的隐式接口（不需要 implements）
  2. EmailSender / SmsSender → 用 struct 实现，注意不需要声明实现了哪个接口
  3. ExecutorService + Future → 用 goroutine + channel 或 sync.WaitGroup 或 errgroup 实现并发控制
  4. 错误收集 → 不要用 panic，用多返回值 + channel 或 slice 收集
  5. 资源管理 → 注意 goroutine 的正确退出和 channel 的关闭