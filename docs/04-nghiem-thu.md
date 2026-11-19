# Backlog và nghiệm thu

Các đường dẫn Go bên dưới là cấu trúc đề xuất cho repo mới. Mỗi dòng là một phần việc có thể review riêng. Lưu bằng chứng theo milestone, gồm SHA gateway, SHA upstream tham khảo, cấu hình đã che bí mật, lệnh chạy, kết quả mong đợi và kết quả thực tế.

## Backlog theo đầu ra

| Việc | Phần bàn giao | Bài nghiệm thu |
|---|---|---|
| G0.1 · Chọn SHA và provider | Hồ sơ baseline, CI và quyết định pilot | Đối chiếu SHA với run CI của base. Lưu link hướng dẫn provider và tên model test. |
| G1.1 · Hợp đồng LiteLLM | `cmd/gateway`, `internal/httpapi`, adapter giả | GET /v1/models trên gateway liệt kê model pilot. Gửi POST /v1/chat/completions qua alias LiteLLM đi đúng gateway. Cả stream và non-stream trả đúng model, kết thúc đúng một lần. Yêu cầu dùng sai key bị từ chối trước dispatch. |
| G1.2 · Cấu hình | `internal/config`, một file mẫu | Nạp model lạ hoặc alias trùng trả lỗi có vị trí. Reload thất bại giữ nguyên cấu hình trước. |
| G2.1 · Đăng nhập và lưu phiên | `internal/oauth`, `internal/store`, một adapter trong `internal/providers` | Dùng callback sai state hoặc device code hết hạn tùy flow. Không tạo session. Đăng nhập hợp lệ và restart vẫn gọi được. Credential được mã hóa khi lưu, khóa mã hóa cấp riêng. |
| G2.2 · Refresh và lỗi 401 | Cùng account chỉ có một lượt refresh tại một thời điểm | 20 request gặp token hết hạn đồng thời tạo một refresh. Kết quả cũ đến muộn không ghi đè token mới. Khi grant bị thu hồi, đánh dấu cần đăng nhập lại đúng account. |
| G3.1 · Đọc quota | `internal/quota` lưu giá trị, đơn vị, scope, nguồn, observed_at, reset_at nếu có | Fixture quota đo được, ước lượng, thiếu và quá hạn cho bốn trạng thái phân biệt được. Quota thiếu không hiển thị bằng 0 hoặc vô hạn. |
| G3.2 · Account eligibility | `internal/accounts`, clock giả, upstream giả có bộ đếm | 429 riêng account A giữ A đến hết Retry-After. 429 thuộc quota nhóm giữ toàn nhóm. Phạm vi mơ hồ thì dừng lượt thử hiện tại và ghi lý do. Hết tài khoản hợp lệ trả 429 và không gửi thêm lượt nào. |
| G4.1 · Quyền vận hành | CLI hoặc config có kiểm tra hợp lệ, event kèm request ID | Tắt model rồi gọi bằng từng alias, đếm 0 dispatch mới. Bật lại có hiệu lực. Trace chỉ chứa ID nội bộ, model, quyết định và thời gian, không chứa token hoặc nội dung chat. |
| G4.2 · Ngắt stream | Hủy upstream khi client ngắt kết nối, đóng stream có kiểm soát | Sau byte SSE đầu tiên, mô phỏng lỗi upstream. Số lần dispatch vẫn là 1. Khi hủy request, tài nguyên được giải phóng và request kế tiếp chạy được. |
| G5.1 · Chuyển tuyến | Hai cấu hình LiteLLM, báo cáo canary và biên bản rollback | Bài kiểm tra trước/sau rollback đi qua đích mong muốn, xác nhận bằng request ID và log đã lọc. |

## Cách chạy và đọc kết quả

Trong repo Go tương lai, nhóm test dùng server HTTP giả và đồng hồ điều khiển được. Chạy `go test ./...` rồi `go test -race ./...`; cả hai phải trả mã 0. Ở G2 chạy thêm smoke test từ **client → LiteLLM → Go → provider thật**. Lưu model, thời gian, mã HTTP, kết quả stream và số lượt refresh; tách rõ kết quả giả lập với bằng chứng live. Đặt retry ở một tầng trong pilot để đo được chính xác số lần gọi.

## Quyết định giữ hay quay lại OmniRoute

1. Trước pilot, đo cùng model bằng cùng tập request tổng hợp qua OmniRoute và Go trong các cửa sổ tương đương. Ghi số mẫu, tỷ lệ 5xx/timeout, p95 thời gian đến token đầu tiên (TTFT), số 429 theo account và lỗi sau khi đã mở stream. Tách lô stream khỏi non-stream.
2. **Ngưỡng pilot đề xuất để Tài duyệt.** Dùng alias riêng trong ít nhất 24 giờ với tối thiểu 100 request mỗi tuyến. Tỷ lệ 5xx/timeout của Go không cao hơn tuyến cũ quá 1 điểm phần trăm, p95 TTFT không cao hơn quá 20%. Yêu cầu 0 lần replay sau khi phát SSE và 0 dispatch vi phạm cooldown. Nếu mẫu ít hoặc provider bị sự cố chung, kéo dài pilot và ghi nhận nguyên nhân.
3. Nếu vi phạm điều kiện hoặc có lỗi phiên lặp lại, trỏ alias về deployment OmniRoute cũ, reload cấu hình LiteLLM và chạy một request thường cùng một request stream. Stream đang chạy được kết thúc theo trạng thái hiện tại. Chỉ tăng traffic sau khi Tài duyệt số đo và bằng chứng rollback.
