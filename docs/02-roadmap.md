# Roadmap từ khảo sát đến chuyển tuyến

**Mốc T0** là ngày G0 đạt. Lịch đề xuất cho một người làm khoảng 8–12 giờ/tuần, có sẵn môi trường LiteLLM và tài khoản thử nghiệm phù hợp. Dự trù 7 tuần từ T0 đến quyết định chuyển tuyến đầu tiên. Thời gian chờ upstream hoặc cấp quyền nằm ngoài ước lượng.

| Mốc | Khi nào · phụ thuộc | Bàn giao | Điều kiện đóng |
|---|---|---|---|
| **G0 · Chốt baseline** | Bắt đầu ngay | SHA nguồn, kết quả base CI, ma trận OAuth, quyết định chọn provider pilot | Queue hết base-red và các check bắt buộc đạt trên cùng SHA. Có cách đăng ký client, tài khoản và đường inference được provider hỗ trợ cho pilot. |
| **G1 · Nối LiteLLM** | Tuần 1 · G0 | Một alias thử nghiệm, gateway Go, upstream giả lập, cấu hình provider/model mẫu | Qua LiteLLM, request thường và SSE đều đúng contract. Cấu hình sai bị từ chối và cấu hình đang chạy vẫn dùng được. |
| **G2 · Một account OAuth** | Tuần 2–3 · G1 | Đăng nhập, lưu phiên, refresh, báo cần đăng nhập lại. Một lệnh chat qua provider pilot | Bài test refresh đồng thời đạt. Restart giữ được phiên. Có bằng chứng một lượt đăng nhập và một request inference thật đã che dữ liệu nhạy cảm. |
| **G3 · Quota và cooldown** | Tuần 4 · G2 | Hiển thị nguồn quota, thời điểm quan sát/reset và lý do tạm dừng, theo đúng account/model hay nhóm quota | 429 có Retry-After chặn dispatch đúng phạm vi. Dữ liệu thiếu hiện unknown. Khi mọi account đều bị chặn, số lần gọi upstream bằng 0. |
| **G4 · Vận hành và MVP** | Tuần 5 · G3 | CLI/config thêm alias, tắt/bật model, xem account, quota và lỗi. Log theo request ID | Tắt model chặn request mới qua mọi alias. Thêm alias cho model đã hỗ trợ không cần build lại. Stream bị đứt sau khi đã phát dữ liệu không bị phát lại. |
| **G5 · Pilot và bàn giao** | Tuần 6–7 · G4 | Báo cáo so sánh với cùng tuyến OmniRoute, cấu hình canary, diễn tập rollback | Canary đạt ngưỡng đã chốt trước khi chạy. Khôi phục được alias cũ và request mới qua OmniRoute thành công. |

**Hai quyết định riêng.** Cuối G4 có MVP dùng được trong staging. Chỉ chuyển alias đang dùng sau G5. Nếu thiếu account hợp lệ hoặc base CI còn đỏ, tiếp tục hoàn thiện ma trận và fixture ở G0, ghi rõ blocker cùng người xử lý.

**Người duyệt đề xuất.** Người triển khai đính kèm test và bằng chứng cho từng mốc. Tài chốt provider ở G0, alias và ngưỡng pilot trước G5, sau đó quyết định giữ gateway Go hay quay về OmniRoute.
