# ephimail
ephimail is an all-in-one disposable email service

## run

```bash
go run cmd/main.go --mail-port 1025 --web-port 8000 --allow-domain localhost
```

## test

```bash
./test-send
```

# todo

- [ ] support for attachments
- [ ] enhance web security posture
- [ ] write smtp test suite
- [ ] implement pub-sub for real time emails
- [ ] hire frontend developer to fix UI
- [ ] implement email TTL
- [ ] TLS on SMTP
- [ ] `reserve` feature, ability to lock a mailbox for a certain amount of time and encrypt incoming emails