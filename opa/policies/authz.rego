package authz

# Mặc định deny tất cả
default allow := false

# Cho phép GET requests
allow if {
    input.method == "GET"
}

# Cho phép authenticated users
allow if {
    input.user.authenticated == true
    input.method != "DELETE"
}

# Admin có thể làm tất cả
allow if {
    input.user.role == "admin"
}

# Chỉ cho phép user thao tác với tài nguyên của chính họ
allow if {
    input.user.authenticated == true
    input.user.id == input.resource.owner_id
    input.method in ["GET", "PUT", "PATCH"]
} 