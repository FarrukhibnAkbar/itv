# ITV

## Loyihani ishga tushirish bo'yicha qo'llanma (O'zbek tilida)

1. Har bir xizmat uchun `.env` faylini yarating:
   - `api_gateway` va `movie_service` xizmatlari uchun kerakli konfiguratsiyalarni `.env` faylida ko'rsating. Bu faylda ma'lumotlar bazasi, portlar va boshqa muhim sozlamalar bo'lishi kerak.

2. Xizmatlarni ishga tushirish:
   - `Makefile` yordamida xizmatlarni ishga tushirish uchun quyidagi buyruqlarni terminalda bajaring:
   ```bash
    make api_gateway
    make movie_service
   ```

3. Loglarni kuzatish:
   - Loyihada loglarni kuzatish uchun maxsus yordamchi funksiya mavjud. Bu funksiya loglarni Telegram kanaliga yuboradi, bu esa tizimdagi xatoliklarni yoki muhim ma'lumotlarni kuzatishni osonlashtiradi.

---

## Project Setup Guide (In English)

1. Create a `.env` file for each service:
   - Provide the necessary configurations in the `.env` file for `api_gateway` and `movie_service`. This file should include database credentials, ports, and other essential settings.

2. Start the services:
   - Use the `Makefile` to start the services by running the following commands in the terminal:
    ```bash
     make api_gateway
     make movie_service
    ```

3. Monitor logs:
   - The project includes a helper function to monitor logs. This function sends logs to a Telegram channel, making it easier to track errors or important events in the system.

---

## Loyihaning tavsifi (O'zbek tilida)

Ushbu loyiha **monorepo arxitekturasi** asosida ishlab chiqilgan bo'lib, bir nechta xizmatlarni bitta repozitoriyada boshqarish imkonini beradi. Loyihaning asosiy maqsadi **Web3** texnologiyalari bilan ishlaydigan, **gRPC** orqali xizmatlararo muloqotni ta'minlaydigan, zamonaviy va kengaytiriladigan tizim yaratishdir.

### Texnologiyalar va arxitektura
- **Monorepo arxitekturasi**: Loyihaning barcha xizmatlari (masalan, `api_gateway`, `movie_service`) bitta repozitoriyada joylashgan. Bu xizmatlararo bog'liqlikni boshqarishni osonlashtiradi va kodni qayta ishlatishni ta'minlaydi.
- **Uber FX**: Loyihada **dependency injection** uchun **Uber FX** ishlatilgan. Bu xizmatlar orasidagi bog'liqlikni avtomatik boshqarish va kodni modular qilish imkonini beradi.
- **HTTP va gRPC**: Tarmoq aloqalari uchun HTTP REST va gRPC ishlatilgan. HTTP REST tashqi mijozlar bilan ishlash uchun, gRPC esa xizmatlararo tezkor va samarali muloqot uchun ishlatiladi.
- **Swagger**: API hujjatlarini avtomatik yaratish va ularga kirishni ta'minlash uchun **Swagger** ishlatilgan.
- **GORM**: Ma'lumotlar bazasi bilan ishlash uchun **GORM** ORM ishlatilgan. Bu ma'lumotlar bazasi operatsiyalarini soddalashtiradi va kodni toza saqlashga yordam beradi.

### Loyihaning asosiy qismlari
1. **API Gateway**:
   - HTTP REST orqali mijozlar bilan muloqot qiladi.
   - Swagger yordamida hujjatlashtirilgan.
   - Xizmatlararo so'rovlarni gRPC orqali `movie_service`ga uzatadi.
2. **Movie Service**:
   - Ma'lumotlar bazasi bilan ishlaydi (GORM orqali).
   - gRPC orqali API Gateway bilan muloqot qiladi.
   - Web3 texnologiyalari bilan integratsiya qilingan.

---

## Project Description (In English)

This project is built using a **monorepo architecture**, which allows managing multiple services within a single repository. The primary goal of the project is to create a modern and scalable system that leverages **Web3** technologies and facilitates inter-service communication via **gRPC**.

### Technologies and Architecture
- **Monorepo Architecture**: All services (e.g., `api_gateway`, `movie_service`) are housed in a single repository. This simplifies inter-service dependencies and promotes code reuse.
- **Uber FX**: Used for **dependency injection**, enabling automatic management of service dependencies and modularizing the codebase.
- **HTTP and gRPC**: Networking is handled using both HTTP REST and gRPC. HTTP REST is used for external client communication, while gRPC ensures fast and efficient inter-service communication.
- **Swagger**: Used for automatic API documentation and providing easy access to API specifications.
- **GORM**: A powerful ORM used for database interactions, simplifying database operations and maintaining clean code.

### Key Components of the Project
1. **API Gateway**:
   - Handles client communication via HTTP REST.
   - Documented using Swagger.
   - Forwards inter-service requests to `movie_service` via gRPC.
2. **Movie Service**:
   - Manages database operations (using GORM).
   - Communicates with the API Gateway via gRPC.
   - Integrated with Web3 technologies.

