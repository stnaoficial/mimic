package com.test.domain.cart;

import com.test.domain.CartBuilder;

public class CartEntity {

    private final String id;

    public CartEntity(CartBuilder builder) {
        this.id = builder.id;
    }

    public String getId() {
        return this.id;
    }

}