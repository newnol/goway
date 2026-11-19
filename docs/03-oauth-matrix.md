# Ma trận tương thích OAuth

Đã đọc mã của năm adapter ở SHA `cc59ad70b2e9`. Tất cả đang ở mức **đối chiếu source**. Quyền dùng client trong gateway riêng, đăng nhập, refresh và inference thực tế còn chờ xác minh theo từng provider.

| Ứng viên | Flow trong OmniRoute | Điểm phải chứng minh trước khi chọn |
|---|---|---|
| Codex | Authorization Code + PKCE, có bước xử lý thông tin account sau đổi token. [Source](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/oauth/providers/codex.ts) | Client và callback dùng được cho ứng dụng riêng. Account đã đăng nhập được phép gọi đúng model pilot. |
| Claude Code | Authorization Code + PKCE, adapter còn nhận dạng thông tin gói khi có dữ liệu. [Source](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/oauth/providers/claude.ts) | Xác minh cơ chế ủy quyền được hỗ trợ cho gateway và xử lý scope thiếu hoặc quyền sử dụng model thay đổi. |
| GitHub Copilot | Device Code, có bước kiểm tra entitlement Copilot sau đăng nhập. [Source](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/oauth/providers/github.ts) | Tách đăng nhập GitHub thành công khỏi quyền dùng Copilot. Kiểm thử account thiếu entitlement và device code hết hạn. |
| Antigravity | Authorization Code, ánh xạ token theo client và tìm project sau khi đổi token. [Source](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/oauth/providers/antigravity.ts) | Xác nhận client được hỗ trợ và cách chọn project. Thiếu project phải trả lỗi chẩn đoán được. |
| Kiro | Device Code qua AWS OIDC, đăng ký client trước khi yêu cầu device code. [Source](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/oauth/providers/kiro.ts) | Xác minh region, đăng ký client, thời hạn client, nhịp poll và việc hủy đăng nhập khi người dùng thoát. |

## Quyết định provider pilot

Ưu tiên khảo sát Codex trước nếu Tài có account phù hợp. Chọn ứng viên đầu tiên có hướng dẫn chính thức cho trường hợp sử dụng này, có cách đăng ký hoặc dùng client hợp lệ và có một model test. Ghi link tài liệu, điều kiện tài khoản, phạm vi quyền và người xác nhận. Nếu cả năm còn thiếu căn cứ, giữ lựa chọn ở trạng thái mở. Mỗi provider bổ sung sau pilot phải qua cùng bộ contract test.

## Hồ sơ source cần chốt ở G0

- Từ registry, lần theo adapter, token refresh, lưu phiên và executor. Ghi file, SHA, hành vi quan sát, điểm chưa rõ và fixture cho từng flow. Kiểm tra lại diff OAuth/quota giữa SHA khảo sát và SHA xanh trước khi chốt. [Registry](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/oauth/providers/index.ts)
- Ghi riêng cơ chế bảo vệ khi hai tiến trình refresh cùng phiên và các nguồn quota upstream phân biệt. Fixture dùng token giả, clock giả và response đã bỏ dữ liệu nhận diện. [Refresh CAS](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/open-sse/services/tokenRefresh/casGuard.ts) · [Quota telemetry](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/src/lib/quota/providerQuotaTelemetry.ts)
