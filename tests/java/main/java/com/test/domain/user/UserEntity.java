package com.test.domain.user;

import com.test.domain.UserBuilder;

public class UserEntity {

    private final String id;

    public UserEntity(UserBuilder builder) {
        this.id = builder.id;
    }

    public String getId() {
        return this.id;
    }

}