# Secure File Vault System

This project implements a secure file vault system with efficient storage, powerful search, and controlled file sharing capabilities. It features a Go backend, PostgreSQL database, and a React.js frontend with TypeScript.

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
  - [Docker Compose Setup](#docker-compose-setup)
- [Database Schema Overview](#database-schema-overview)
- [API Schema](#api-schema)
- [Design/Architecture Writeup](#designarchitecture-writeup)
- [Bonus Features](#bonus-features)
- [UAT Checklist and Automated Tests](#uat-checklist-and-automated-tests)

## Features

The application implements the following core features:

1.  **File Deduplication**: Detects duplicate uploads using SHA-256 content hashing. Stores references instead of duplicate content and displays per-user storage savings.
2.  **File Uploads**: Supports single and multiple file uploads with drag-and-drop functionality in the frontend. Validates file content against declared MIME types.
3.  **File Management & Sharing**:
    *   **List and View Files**: Users can view their owned files with detailed metadata (uploader, size, upload date, deduplication info).
    *   **Organize Files into Folders**: Users can organize files into folders.
    *   **File and Folder Sharing**:
        *   Public sharing with a link.
        *   Private files/folders visible only to the owner.
        *   Sharing with specific users (bonus).
    *   **Public File Statistics**: Publicly shared files show download counters.
    *   **Delete Files**: Only the uploader can delete their files. Deduplicated files respect reference counts.
4.  **Search & Filtering**: Search by filename and filter by MIME type, size range, date range, tags, and uploader's name. Multiple filters can be combined. Queries are optimized for performance.
5.  **Rate Limiting & Quotas**: Enforces per-user API rate limits (2 calls/second, configurable) and storage quotas (10 MB, configurable). Returns proper error codes/messages.
6.  **Storage Statistics**: Displays total storage used (deduplicated), original storage usage, and storage savings in bytes and percentage.
7.  **Admin Panel**: Admins can upload files, share with other users, list all files with uploader details, and view download counts and usage stats.

## Tech Stack

*   **Backend**: Go (Golang)
*   **API Layer**: REST
*   **Database**: Supabase (PostgreSQL)
*   **File Storage**: Supabase Storage (S3 Bucket)
*   **Frontend**: React.js with TypeScript
*   **Containerization**: Docker Compose

## Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go             # Backend entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers.go         # API handlers for various routes
│   │   │   ├── middleware.go       # Middleware for authentication, rate limiting, etc.
│   │   │   └── routes.go           # API route definitions
│   │   ├── database/
│   │   │   └── database.go         # Database connection and utility functions
│   │   ├── email/
│   │   │   └── email.go            # Email sending functionalities
│   │   ├── handlers/
│   │   │   ├── admin.go            # Admin specific handlers
│   │   │   ├── files.go            # File related handlers (upload, download, delete, share)
│   │   │   └── users.go            # User related handlers (auth, profile)
│   │   └── models/
│   │       ├── file.go             # File and related data models
│   │       └── user.go             # User and related data models
│   ├── go.mod
│   ├── go.sum
│   └── .env                        # Environment variables
├── frontend/
│   ├── public/
│   │   ├── logo.png
│   │   └── vite.svg
│   ├── src/
│   │   ├── assets/
│   │   │   └── react.svg
│   │   ├── components/             # React components
│   │   │   ├── auth/               # Authentication related components
│   │   │   │   ├── GuestRoute.tsx
│   │   │   │   └── ProtectedRoute.tsx
│   │   │   ├── AllFiles.tsx
│   │   │   ├── EmailVerificationPage.tsx
│   │   │   ├── FileDetailSidebar.tsx
│   │   │   ├── FilterPopover.tsx
│   │   │   ├── HomePage.tsx
│   │   │   ├── Layout.tsx
│   │   │   ├── LoginPage.tsx
│   │   │   ├── ProfilePopover.tsx
│   │   │   ├── PubliclySharedFiles.tsx
│   │   │   ├── PublicShareView.tsx
│   │   │   ├── SettingsPage.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── SignUpPage.tsx
│   │   ├── contexts/               # React contexts
│   │   │   └── ThemeContext.tsx
│   │   ├── hooks/                  # Custom React hooks
│   │   │   ├── useClickOutside.tsx
│   │   │   ├── useConfirmationDialog.tsx
│   │   │   └── useFileUpload.tsx
│   │   ├── utils/                  # Utility functions
│   │   │   └── formatBytes.ts
│   │   ├── App.tsx                 # Main React application component
│   │   ├── index.css               # Global CSS
│   │   ├── main.tsx                # Frontend entry point
│   │   └── vite-env.d.ts
│   ├── .gitignore
│   ├── eslint.config.js
│   ├── index.html
│   ├── package.json
│   ├── package-lock.json
│   ├── postcss.config.js
│   ├── tailwind.config.js
│   ├── tsconfig.app.json
│   ├── tsconfig.json
│   ├── tsconfig.node.json
│   └── vite.config.ts
├── .gitignore
├── package-lock.json
├── README.md                       # This documentation file
├── test.txt
```

## Setup Instructions

### Prerequisites

*   Docker and Docker Compose
*   Go (Golang) 1.18+
*   Node.js 18+ and npm

### Backend Setup

1.  Navigate to the `backend` directory:
    ```bash
    cd backend
    ```
2.  Create a `.env` file in the `backend` directory and populate it with the following environment variables:
    ```
    DATABASE_URL="postgresql://postgres:zocw5Yo44rsghwrN@ounyoiaygoidjgcwpxpv.supabase.co:5432/postgres"
    SUPABASE_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im91bnlvaWF5Z29pZGpnY3dweHB2Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc1ODA0MDIxOSwiZXhwIjoyMDczNjE2MjE5fQ.OZF3fsdZni6uq-PLnhShqEjmjHKTjatA3p0aLiq4WOs"
    SUPABASE_URL="https://ounyoiaygoidjgcwpxpv.supabase.co"
    SUPABASE_REST_URL="https://ounyoiaygoidjgcwpxpv.supabase.co/rest/v1"
    SUPABASE_STORAGE_URL="https://ounyoiaygoidjgcwpxpv.storage.supabase.co/storage/v1"
    SUPABASE_S3_ACCESS_KEY="31d493b5537d1682d0364f913d5c2df7"

    # Email Service Configuration
    SMTP_HOST=smtp.gmail.com
    SMTP_USER="ramkumar.arcot2022@vitstudent.ac.in"
    SMTP_PASS="odks vmhe nqjo zyil"
    ```
    **Note**: For security, replace `SUPABASE_KEY`, `SMTP_USER`, and `SMTP_PASS` with your actual credentials.
3.  Install Go dependencies:
    ```bash
    go mod tidy
    ```
4.  Run database migrations (if not using Docker Compose for the database):
    ```bash
    # Example using a migration tool like 'migrate'
    migrate -path database/migrations -database "postgres://user:password@host:port/dbname?sslmode=disable" up
    ```
5.  Run the backend server:
    ```bash
    go run cmd/server/main.go
    ```

### Frontend Setup

1.  Navigate to the `frontend` directory:
    ```bash
    cd frontend
    ```
2.  Install Node.js dependencies:
    ```bash
    npm install
    ```
3.  Start the frontend development server:
    ```bash
    npm run dev
    ```
    The frontend application will typically be available at `http://localhost:5173`.

### Docker Compose Setup

For a fully containerized local development environment:

### Sample Credentials

For testing purposes, you can use the following sample credentials:

*   **Test User**:
    *   Username: `user_one`
    *   Email: `user_one@test.com`
    *   Password: `userone`

*   **Admin User**:
    *   Username: `admin`
    *   Email: `admin@test.com`
    *   Password: `admin`

1.  Ensure Docker and Docker Compose are installed and running.
2.  From the project root directory, build and start the services:
    ```bash
    docker-compose up --build
    ```
    This will start the PostgreSQL database, backend API, and frontend application.

## Database Schema Overview

The PostgreSQL database schema, hosted on Supabase, is designed to support file storage, deduplication, user management, sharing, and analytics. File contents are stored in Supabase Storage (S3 bucket).

```sql
-- Users Table: Stores user information, including authentication details, quotas, and admin status.
CREATE TABLE public.users (
  user_id uuid NOT NULL DEFAULT gen_random_uuid(),
  username character varying NOT NULL UNIQUE,
  email character varying NOT NULL UNIQUE,
  password_hash text NOT NULL,
  is_admin boolean DEFAULT false,
  storage_quota bigint DEFAULT 10485760, -- Default 10MB
  created_at timestamp without time zone DEFAULT now(),
  first_name character varying,
  last_name character varying,
  date_of_birth date,
  phone_number character varying,
  last_login timestamp without time zone,
  email_verified boolean DEFAULT false,
  phone_verified boolean DEFAULT false,
  status character varying DEFAULT 'active'::character varying,
  otp character varying,
  otp_expires_at timestamp without time zone,
  rate_limit integer DEFAULT 2, -- Default 2 calls per second
  CONSTRAINT users_pkey PRIMARY KEY (user_id)
);

-- FileContents Table: Stores unique file content information, enabling deduplication.
CREATE TABLE public.file_contents (
  content_id uuid NOT NULL DEFAULT gen_random_uuid(),
  hash_sha256 character NOT NULL UNIQUE, -- SHA-256 hash of the file content
  size bigint NOT NULL,
  mime_type character varying,
  storage_path text NOT NULL, -- Path to the actual stored file
  reference_count integer NOT NULL DEFAULT 0, -- Number of files referencing this content
  created_at timestamp without time zone DEFAULT now(),
  CONSTRAINT file_contents_pkey PRIMARY KEY (content_id)
);

-- Files Table: Represents a user's logical file entry, pointing to a unique file_content.
CREATE TABLE public.files (
  file_id uuid NOT NULL DEFAULT gen_random_uuid(),
  owner_id uuid,
  content_id uuid,
  filename character varying NOT NULL,
  is_deleted boolean DEFAULT false,
  created_at timestamp without time zone DEFAULT now(),
  CONSTRAINT files_pkey PRIMARY KEY (file_id),
  CONSTRAINT files_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(user_id),
  CONSTRAINT files_content_id_fkey FOREIGN KEY (content_id) REFERENCES public.file_contents(content_id)
);

-- Folders Table: Allows users to organize their files into a hierarchical structure.
CREATE TABLE public.folders (
  folder_id uuid NOT NULL DEFAULT gen_random_uuid(),
  owner_id uuid,
  parent_id uuid, -- Self-referencing for nested folders
  name character varying NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  CONSTRAINT folders_pkey PRIMARY KEY (folder_id),
  CONSTRAINT folders_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(user_id),
  CONSTRAINT folders_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.folders(folder_id)
);

-- FileFolderMapping Table: Links files to folders.
CREATE TABLE public.file_folder_mapping (
  file_id uuid NOT NULL,
  folder_id uuid NOT NULL,
  CONSTRAINT file_folder_mapping_pkey PRIMARY KEY (folder_id, file_id),
  CONSTRAINT file_folder_mapping_file_id_fkey FOREIGN KEY (file_id) REFERENCES public.files(file_id),
  CONSTRAINT file_folder_mapping_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.folders(folder_id)
);

-- Tags Table: Stores predefined or user-defined tags for files.
CREATE TABLE public.tags (
  tag_id uuid NOT NULL DEFAULT gen_random_uuid(),
  name character varying NOT NULL UNIQUE,
  CONSTRAINT tags_pkey PRIMARY KEY (tag_id)
);

-- FileTags Table: Many-to-many relationship between files and tags.
CREATE TABLE public.file_tags (
  file_id uuid NOT NULL,
  tag_id uuid NOT NULL,
  CONSTRAINT file_tags_pkey PRIMARY KEY (file_id, tag_id),
  CONSTRAINT file_tags_file_id_fkey FOREIGN KEY (file_id) REFERENCES public.files(file_id),
  CONSTRAINT file_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(tag_id)
);

-- Shares Table: Manages file and folder sharing, including public links and specific user shares.
CREATE TABLE public.shares (
  share_id uuid NOT NULL DEFAULT gen_random_uuid(),
  file_id uuid,
  folder_id uuid,
  is_public boolean DEFAULT false,
  shared_with uuid, -- Null for public shares, references user_id for specific shares
  created_at timestamp without time zone DEFAULT now(),
  share_token text, -- Unique token for public access
  download_count integer NOT NULL DEFAULT 0,
  CONSTRAINT shares_pkey PRIMARY KEY (share_id),
  CONSTRAINT shares_file_id_fkey FOREIGN KEY (file_id) REFERENCES public.files(file_id),
  CONSTRAINT shares_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.folders(folder_id),
  CONSTRAINT shares_shared_with_fkey FOREIGN KEY (shared_with) REFERENCES public.users(user_id)
);

-- DownloadLogs Table: Records each download event for public files.
CREATE TABLE public.download_logs (
  log_id bigint NOT NULL DEFAULT nextval('download_logs_log_id_seq'::regclass),
  file_id uuid,
  downloader_id uuid, -- Null for anonymous public downloads
  downloaded_at timestamp without time zone DEFAULT now(),
  CONSTRAINT download_logs_pkey PRIMARY KEY (log_id),
  CONSTRAINT download_logs_file_id_fkey FOREIGN KEY (file_id) REFERENCES public.files(file_id),
  CONSTRAINT download_logs_downloader_id_fkey FOREIGN KEY (downloader_id) REFERENCES public.users(user_id)
);

-- APIUsage Table: Logs API calls for rate limiting and analytics.
CREATE TABLE public.api_usage (
  usage_id bigint NOT NULL DEFAULT nextval('api_usage_usage_id_seq'::regclass),
  user_id uuid,
  endpoint character varying NOT NULL,
  called_at timestamp without time zone DEFAULT now(),
  CONSTRAINT api_usage_pkey PRIMARY KEY (usage_id),
  CONSTRAINT api_usage_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id)
);
```

## API Schema

The backend exposes a RESTful API. Below are some key endpoints and their functionalities. For a complete and interactive API documentation, please refer to the OpenAPI (Swagger) specification generated by the backend or use a tool like Postman with the provided collection.

**Base URL**: `/api/v1` (configurable)

### Authentication

*   `POST /register`: Register a new user.
    *   **Request Body**: `{ "username": "...", "email": "...", "password": "..." }`
    *   **Response**: `{ "message": "User registered successfully" }`
*   `POST /login`: Authenticate a user and receive a JWT token.
    *   **Request Body**: `{ "email": "...", "password": "..." }`
    *   **Response**: `{ "token": "...", "user": { ... } }`
*   `POST /verify-email`: Verify user email using OTP.
    *   **Request Body**: `{ "email": "...", "otp": "..." }`
    *   **Response**: `{ "message": "Email verified successfully" }`

### File Management

*   `POST /files/upload`: Upload one or more files.
    *   **Request Body**: `multipart/form-data` with file(s)
    *   **Response**: `{ "message": "Files uploaded successfully", "files": [...] }`
*   `GET /files`: List all files owned by the authenticated user.
    *   **Query Parameters**: `filename`, `mime_type`, `min_size`, `max_size`, `start_date`, `end_date`, `tags`, `uploader_name` for filtering.
    *   **Response**: `[ { "file_id": "...", "filename": "...", "size": "...", ... } ]`
*   `GET /files/{file_id}`: Get details of a specific file.
    *   **Response**: `{ "file_id": "...", "filename": "...", "size": "...", ... }`
*   `DELETE /files/{file_id}`: Delete a file. Only the owner can delete. Deduplication reference counts are handled.
    *   **Response**: `{ "message": "File deleted successfully" }`
*   `GET /files/{file_id}/download`: Download a specific file.
    *   **Response**: File content.

### Folder Management

*   `POST /folders`: Create a new folder.
    *   **Request Body**: `{ "name": "...", "parent_id": "..." (optional) }`
    *   **Response**: `{ "folder_id": "...", "name": "...", ... }`
*   `GET /folders`: List all folders owned by the authenticated user.
    *   **Response**: `[ { "folder_id": "...", "name": "...", ... } ]`
*   `POST /folders/{folder_id}/add-file/{file_id}`: Add an existing file to a folder.
    *   **Response**: `{ "message": "File added to folder successfully" }`

### Sharing

*   `POST /shares/file/{file_id}/public`: Create a public share link for a file.
    *   **Response**: `{ "share_link": "...", "share_token": "..." }`
*   `POST /shares/folder/{folder_id}/public`: Create a public share link for a folder.
    *   **Response**: `{ "share_link": "...", "share_token": "..." }`
*   `GET /public/share/{share_token}`: Access publicly shared file/folder content.
    *   **Response**: File/folder details and content.
*   `GET /public/share/{share_token}/download`: Download a publicly shared file.
    *   **Response**: File content.

### Statistics

*   `GET /stats/storage`: Get user storage statistics (total, original, savings).
    *   **Response**: `{ "total_used": ..., "original_used": ..., "savings_bytes": ..., "savings_percentage": ... }`
*   `GET /stats/public-downloads/{file_id}`: Get download count for a public file.
    *   **Response**: `{ "download_count": ... }`

### Admin Endpoints (Requires Admin Role)

*   `GET /admin/users`: List all users.
*   `GET /admin/files`: List all files across all users.
*   `POST /admin/files/upload-and-share`: Admin uploads a file and shares it with a specific user.
    *   **Request Body**: `multipart/form-data` with file, `shared_with_user_id`
*   `GET /admin/stats/usage`: View overall system usage statistics.

## Design/Architecture Writeup

The application follows a layered architecture, separating concerns into distinct components:

*   **Frontend (React.js)**: A responsive Single Page Application (SPA) built with React and TypeScript. It interacts with the backend API to provide a user-friendly interface for file management, uploads, sharing, and search. State management is handled using React Context and local component state. Custom hooks are used for reusable logic.
*   **Backend (Go)**: A RESTful API built with Go, designed for performance and scalability.
    *   **`cmd/server/main.go`**: The entry point, responsible for initializing the server, database connection, and routing.
    *   **`internal/api`**: Contains HTTP route definitions (`routes.go`), request handlers (`handlers.go`), and middleware (`middleware.go`) for authentication, authorization, and rate limiting.
    *   **`internal/handlers`**: Business logic for specific domains (users, files, admin). These handlers interact with the database and other internal services.
    *   **`internal/database`**: Manages PostgreSQL database connections, queries, and transactions. Uses `pgx` for efficient database interaction.
    *   **`internal/models`**: Defines the data structures (structs) that map to database tables and are used throughout the application.
    *   **`internal/email`**: Handles sending emails, e.g., for OTP verification.
*   **Database (PostgreSQL)**: A robust relational database used for persistent storage. The schema is designed to support deduplication (via `file_contents` and `files` tables), hierarchical folder structures, and detailed logging for downloads and API usage.
*   **Deduplication Logic**: When a file is uploaded, its SHA-256 hash is calculated. The `file_contents` table is checked for an existing entry with the same hash. If found, a new `files` entry is created referencing the existing `content_id`, and the `reference_count` in `file_contents` is incremented. If not found, the file content is stored, a new `file_contents` entry is created, and then a `files` entry references it. Deletion decrements the `reference_count`, and the actual content is only removed when `reference_count` reaches zero.
*   **Rate Limiting**: Implemented as middleware, tracking API calls per user within a time window using an in-memory store or a distributed cache (e.g., Redis) for production.
*   **Storage Quotas**: Enforced during file uploads by checking the user's current storage against their `storage_quota` defined in the `users` table.
*   **Security**: JWT-based authentication, password hashing (bcrypt), MIME type validation, and access control for file operations and admin functionalities.

## Bonus Features

*   **Sharing with specific users**: The `shares` table includes a `shared_with` column, allowing files/folders to be shared with specific user IDs.
*   **Admin Panel with graphs, sharing and analytics**: The admin handlers and frontend components provide basic functionalities for admin users to manage files, users, and view usage statistics. Further enhancements for graphs and detailed analytics can be built upon this foundation.
*   **File Previews**: While not fully implemented, the `mime_type` stored in `file_contents` provides the necessary information to implement file previews for supported types (e.g., images, PDFs) in the frontend.
*   **Role-Based Access Control (RBAC)**: Implemented through the `is_admin` flag in the `users` table and checked in backend middleware and handlers.
*   **Audit Logs**: `download_logs` and `api_usage` tables serve as a foundation for audit logging of key activities.

## UAT Checklist and Automated Tests

A comprehensive UAT (User Acceptance Testing) checklist would cover all core features and bonus functionalities from both user and admin perspectives.

**Example UAT Checklist Items:**

*   **User Registration & Login**:
    *   [ ] Register a new user successfully.
    *   [ ] Log in with valid credentials.
    *   [ ] Attempt to log in with invalid credentials (expect error).
    *   [ ] Verify email via OTP.
*   **File Uploads**:
    *   [ ] Upload a single file.
    *   [ ] Upload multiple files.
    *   [ ] Drag-and-drop a file for upload.
    *   [ ] Upload a duplicate file (expect deduplication).
    *   [ ] Attempt to upload a file exceeding storage quota (expect error).
    *   [ ] Upload a file with mismatched MIME type (e.g., `.txt` renamed to `.jpg`, expect error).
*   **File Management**:
    *   [ ] View a list of owned files.
    *   [ ] View detailed metadata for a file.
    *   [ ] Create a new folder.
    *   [ ] Move a file into a folder.
    *   [ ] Delete an owned file (expect success).
    *   [ ] Attempt to delete a file owned by another user (expect error).
    *   [ ] Delete a deduplicated file (verify `reference_count` decrement).
    *   [ ] Download an owned file.
*   **File Sharing**:
    *   [ ] Create a public share link for a file.
    *   [ ] Access a publicly shared file via its link (unauthenticated).
    *   [ ] Download a publicly shared file (unauthenticated).
    *   [ ] Verify download count increments for public files.
    *   [ ] Create a public share link for a folder.
    *   [ ] Access a publicly shared folder via its link (unauthenticated).
    *   [ ] Share a file with a specific user (authenticated).
    *   [ ] Verify the specific user can access the shared file.
*   **Search & Filtering**:
    *   [ ] Search files by filename.
    *   [ ] Filter files by MIME type.
    *   [ ] Filter files by size range.
    *   [ ] Filter files by date range.
    *   [ ] Filter files by tags.
    *   [ ] Combine multiple filters (e.g., filename + MIME type + date range).
*   **Rate Limiting & Quotas**:
    *   [ ] Exceed API rate limit (expect 429 Too Many Requests).
    *   [ ] Exceed storage quota during upload (expect error).
*   **Storage Statistics**:
    *   [ ] View total storage used.
    *   [ ] View original storage usage.
    *   [ ] View storage savings in bytes and percentage.
*   **Admin Panel**:
    *   [ ] Log in as an admin user.
    *   [ ] View all users.
    *   [ ] View all files (including uploader details).
    *   [ ] Upload a file as admin and share with another user.
    *   [ ] View download counts and usage stats from admin panel.

**Automated Tests**:

Automated tests can be implemented at various levels:

*   **Unit Tests (Go)**: Using Go's built-in `testing` package for individual functions and methods in the backend.
*   **Integration Tests (Go)**: Testing the interaction between backend components (e.g., handlers with database).
*   **API Tests (Postman/Newman)**: All API endpoints have been thoroughly tested using Postman. A Postman collection can be created to test all API endpoints, including authentication, file operations, sharing, and admin functionalities. Newman (Postman's CLI runner) can automate these tests in a CI pipeline.
*   **E2E Tests (Playwright/Cypress)**: For the frontend, Playwright or Cypress can simulate user interactions to test critical user flows (e.g., user registration -> login -> file upload -> file sharing -> file deletion). These tests would ensure the frontend and backend integrate correctly.

**CI Pipeline Integration**:

A GitHub Actions workflow can be configured to:

1.  Run backend Go unit and integration tests on push/pull request.
2.  Run frontend Node.js unit tests (e.g., Jest/React Testing Library) on push/pull request.
3.  Execute API tests using Newman against a deployed staging environment or a locally spun-up Docker Compose environment.
4.  Run Playwright/Cypress E2E tests against the deployed frontend.

This setup ensures that code changes are continuously validated, maintaining code quality and preventing regressions.

[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/2xw7QaEj)

## Screenshots

Here are some screenshots of the application:

![Homepage](frontend/src/assets/homepage.png)
![All Files Page](frontend/src/assets/allFilesPage.png)
![Settings Page](frontend/src/assets/settingsPage.png)
![Shared Publicly Page](frontend/src/assets/sharedPubliclyPage.png)
