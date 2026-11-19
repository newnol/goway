# G0 · Baseline OmniRoute và provider pilot

**Trạng thái ngày 19/11/2026 · Đang chờ.** Đã đối chiếu mã nguồn, CI và luồng OAuth chính thức. Còn thiếu lượt kiểm tra base chạy đủ trên đúng SHA và xác nhận quyền sử dụng ChatGPT plan của tài khoản pilot.

## Baseline OmniRoute

| Mục | Bằng chứng |
|---|---|
| Đầu nhánh khảo sát | `release/v3.8.55`, SHA `adca22a0395ef2bc41acba5610a1532eda621b22`, phiên bản `3.8.55`. [Commit](https://github.com/cheaper-inference/OmniRoute/commit/adca22a0395ef2bc41acba5610a1532eda621b22) · [package.json](https://github.com/cheaper-inference/OmniRoute/blob/adca22a0395ef2bc41acba5610a1532eda621b22/package.json) |
| Tình trạng queue | Queue báo chạy lại lúc 17:23 giờ Việt Nam khi base chuyển từ `cc59ad70b2e9` sang `b728fe2de5e7`. Run đi kèm là lượt `Tick` của queue. [Queue #20583](https://github.com/cheaper-inference/OmniRoute/issues/20583) · [Queue tick](https://github.com/cheaper-inference/OmniRoute/actions/runs/44505143879) |
| Kiểm tra trên `adca22a0395e` | Hai workflow CI và Quality Gates báo thành công tổng thể, nhưng danh sách job cho thấy Unit Tests, Build và hai job Gate bị bỏ qua. Chưa có bằng chứng để đánh dấu SHA này đạt điều kiện G0. [CI](https://github.com/cheaper-inference/OmniRoute/actions/runs/44505308453) · [Quality Gates](https://github.com/cheaper-inference/OmniRoute/actions/runs/44505308372) |
| Thay đổi liên quan | Diff từ SHA khảo sát cũ sang SHA mới không có thay đổi trong `src/lib/oauth`, `src/lib/quota`, `open-sse/services/tokenRefresh`, `open-sse/executors` và `src/lib/db`. Có thể tiếp tục dùng [ma trận 5 adapter](03-oauth-matrix.md) làm tài liệu tham khảo. [So sánh hai SHA](https://github.com/cheaper-inference/OmniRoute/compare/cc59ad70b2e9a343f26860e0e040281cb433a37a...adca22a0395ef2bc41acba5610a1532eda621b22) |

## Phương án pilot

**Ưu tiên Sign in with ChatGPT cho một tài khoản trên GoWay tự chạy cục bộ.** OpenAI có luồng đăng ký client cho ứng dụng mã nguồn mở, dùng Authorization Code với PKCE. Tài khoản cần cấp quyền `chatgpt.tokens.use.direct`. Sau khi nhận token, GoWay lấy model catalog qua `GET /v1/models` và gửi inference qua `POST /v1/responses` trên `api.openai.com`. [Đăng ký và đăng nhập](https://developers.openai.com/siwc/token-sharing-open-source/sign-in) · [Model và inference](https://developers.openai.com/siwc/token-sharing-open-source/models-and-inference)

- **Model dự kiến.** Dùng `gpt-6.1-sol` nếu catalog của tài khoản pilot liệt kê model này. Nếu thiếu, ghi lại slug của một model text có trong catalog trước khi chốt. Chưa xác minh quyền tài khoản hoặc gọi inference thật.
- **Ảnh hưởng tới gateway.** Upstream pilot dùng Responses SSE với `stream: true` và `store: false`. GoWay cần chuyển event thành Chat Completions SSE cho LiteLLM, hoặc gom text thành JSON khi client gọi ở chế độ thường. Chỉ trả kết thúc thành công khi nhận `response.completed`. [Giới hạn bản preview](https://developers.openai.com/siwc/token-sharing-open-source/preview-limitations)
- **Phương án dự phòng.** GitHub Copilot SDK có luồng GitHub OAuth App hoặc GitHub App được tài liệu hóa. Cần kiểm tra subscription và chốt riêng nếu tài khoản ChatGPT không đủ điều kiện. [GitHub OAuth setup](https://docs.github.com/en/copilot/how-tos/copilot-sdk/setup/github-oauth)

## Tiêu chí kiểm thử cho adapter pilot

| Tình huống mô phỏng | Kết quả mong đợi |
|---|---|
| Callback sai `state` hoặc người dùng từ chối | Không đổi code lấy token, phiên hiện tại vẫn dùng được. |
| Đăng nhập thành công nhưng thiếu quyền dùng ChatGPT plan | Không gửi inference. Báo rõ thiếu quyền. |
| Model dự kiến vắng mặt trong catalog | Không dispatch theo tên mặc định. Chọn model từ catalog trước khi thử tiếp. |
| Stream phát text rồi trả `response.failed` hoặc kết nối đóng trước `response.completed` | Báo lỗi cho client, không trả dấu kết thúc thành công và không tự phát lại stream. |
| Refresh trả `invalid_grant` | Đánh dấu đúng tài khoản cần đăng nhập lại. |

## Việc còn lại để đóng G0

1. **Người triển khai** ghi nhận lượt base CI chạy đủ các check bắt buộc trên cùng SHA, kèm đường dẫn run. Nếu đầu nhánh đổi, kiểm tra lại diff OAuth và quota.
2. **Tài** xác nhận tài khoản pilot và quyền dùng ChatGPT plan. Chốt một model text từ catalog. Lượt đăng nhập và inference qua adapter GoWay sẽ được nghiệm thu ở G2.
