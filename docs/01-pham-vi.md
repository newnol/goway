# Phạm vi và điểm xuất phát

**Đích đến.** Tài có thể chuyển một model alias đang dùng trong LiteLLM sang gateway Go, đăng nhập lại khi phiên OAuth hết hiệu lực, thấy tài khoản nào đang chạm quota và quay về OmniRoute bằng thay đổi cấu hình LiteLLM.

## Ai chịu trách nhiệm phần nào?

| Thành phần | Trách nhiệm trong bản đầu |
|---|---|
| LiteLLM | Giữ API cho ứng dụng, alias và chính sách chọn model/provider. Trỏ một deployment thử nghiệm sang gateway Go. |
| Gateway Go | Nhận request cho một model đã ánh xạ, chọn account đủ điều kiện trong đúng provider, quản lý OAuth, quota, cooldown và trả kết quả thường hoặc stream. |
| Provider adapter | Khai báo model hỗ trợ, flow đăng nhập, cách refresh, gọi inference và đọc tín hiệu quota/lỗi đặc thù. Thêm alias hoặc tắt model bằng cấu hình. Thay giao thức mới cần adapter Go và contract test. |
| OmniRoute hiện tại | Tiếp tục phục vụ các tuyến còn lại và là đích rollback của alias thử nghiệm. |

## Trạng thái đối chiếu ngày 19/11/2026

Các dòng dưới ghi lại thời điểm lập kế hoạch. Xem [lần kiểm tra G0 sau đó](00-g0-baseline.md) để biết SHA mới, tình trạng CI và phương án pilot.

- **Upstream dùng để đọc source** là nhánh `release/v3.8.55`, phiên bản trong `package.json` là `3.8.55`, SHA `cc59ad70b2e9a343f26860e0e040281cb433a37a`. [Commit](https://github.com/cheaper-inference/OmniRoute/commit/cc59ad70b2e9a343f26860e0e040281cb433a37a) · [package.json](https://github.com/cheaper-inference/OmniRoute/blob/cc59ad70b2e9a343f26860e0e040281cb433a37a/package.json)
- **G0 đang chờ CI.** Queue báo dừng lúc 16:16 giờ Việt Nam do `base-red` trên chính SHA này. Lần kiểm tra base ghi nhận lỗi ở Unit Tests (6/8). Khi đầu nhánh đổi, lưu SHA mới và kết quả các check bắt buộc trước khi chọn làm baseline. [Queue #20583](https://github.com/cheaper-inference/OmniRoute/issues/20583) · [Base check](https://github.com/cheaper-inference/OmniRoute/actions/runs/44494913411)
- **Fork TypeScript.** PR #2 (audit) và #3 (quy tắc lifecycle) đã merge. PR #4 (lưu policy bằng SQLite và kiểm soát cập nhật đồng thời) đang mở. Có thể tái dùng các tình huống kiểm thử của fork làm đầu vào thiết kế Go. [#2](https://github.com/newnol/OmniRoute/pull/2) · [#3](https://github.com/newnol/OmniRoute/pull/3) · [#4](https://github.com/newnol/OmniRoute/pull/4)

## Chọn phần nào từ kế hoạch fork v3?

Giữ quy tắc tắt model ở lần kiểm tra cuối trước khi gửi, giới hạn phạm vi ảnh hưởng của 401/429 và ghi rõ quota đến từ đâu. Lùi dashboard quản trị, tối ưu catalog 10k model, ngân sách tiền theo ví dùng chung, multi-replica và agent phân tích sang sau pilot. Cắt MVP ở một provider, một model chat, một instance và một writer. Chỉ mở rộng khi provider thứ hai cho thấy phần nào đáng dùng chung.
