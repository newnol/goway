# Gateway Go phía sau LiteLLM

Bộ kế hoạch ngày 19/11/2026 cho dự án Omniroute Fork. Mục tiêu là thay từng tuyến OAuth của OmniRoute bằng một gateway Go nhỏ, có thể vận hành và quay lui theo từng model alias trong LiteLLM.

| Khi cần | Đọc |
|---|---|
| Xem kết quả G0 và điều kiện còn thiếu | [Baseline và provider pilot](00-g0-baseline.md) |
| Xác định phạm vi, trạng thái upstream và phần kế thừa từ fork | [Phạm vi và điểm xuất phát](01-pham-vi.md) |
| Chốt lịch, phụ thuộc và người duyệt | [Roadmap G0–G5](02-roadmap.md) |
| Chọn provider thử nghiệm và biết cần đọc gì trong source | [Ma trận tương thích OAuth](03-oauth-matrix.md) |
| Chia việc, chạy nghiệm thu và quyết định chuyển traffic | [Backlog và nghiệm thu](04-nghiem-thu.md) |

Mã G dùng cho gateway Go. M0–M4 trong bộ kế hoạch v3 thuộc fork TypeScript. Mỗi lần cập nhật tiến độ, ghi ngày, SHA và đường dẫn bằng chứng cạnh milestone tương ứng. Ngày trong roadmap là ước lượng tương đối.

[Bản PDF ngày 19/11/2026](OmniRoute-Go-Roadmap.pdf) · [Chạy thử GoWay](../README.md#quick-start)
