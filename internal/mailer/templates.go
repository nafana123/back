package mail

func buildVerifyEmailHTML(code string) string {
	return `<!DOCTYPE html>
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
</head>
<body style="margin:0;padding:0;font-family:Arial,Helvetica,sans-serif;background:#ffffff;">

  <div style="max-width:700px;margin:40px auto;padding:40px 24px;background:#ffffff;border-radius:14px;box-shadow:0 10px 30px rgba(0,0,0,0.25);text-align:center;">

    <div style="font-size:14px;letter-spacing:3px;color:#888;margin-bottom:10px;">
      TOURNIFY
    </div>

    <h1 style="margin:0 0 18px 0;font-size:28px;font-weight:800;color:#111;">
      Подтверждение регистрации
    </h1>

    <p style="margin:0 0 25px 0;font-size:16px;color:#444;line-height:1.6;">
      Чтобы завершить регистрацию, используйте одноразовый код ниже.
    </p>

    <div style="
      display:inline-block;
      margin:20px 0;
      padding:14px 28px;
      font-size:34px;
      font-weight:800;
      letter-spacing:6px;
      color:#000000;
      background:#f5f5f5;
      border-radius:12px;
      border:1px solid #e5e5e5;
    ">
      ` + code + `
    </div>

    <p style="margin:18px 0 0 0;font-size:14px;color:#666;">
      Код действителен в течении <b>5 минут</b>
    </p>

    <p style="margin:25px 0 0 0;font-size:13px;color:#777;line-height:1.5;">
      Если вы не создавали аккаунт в Tournify, просто проигнорируйте это письмо.
    </p>

    <div style="margin-top:40px;font-size:12px;color:#aaa;">
      © 2026 Tournify • Gaming tournaments platform
    </div>
  </div>

</body>
</html>`
}
