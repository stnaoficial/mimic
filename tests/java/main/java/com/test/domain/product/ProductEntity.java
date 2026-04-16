package com.test.domain.product;

import com.test.domain.ProductBuilder;

public class ProductEntity {

    private final String id;

    public ProductEntity(ProductBuilder builder) {
        this.id = builder.id;
    }

    public String getId() {
        return this.id;
    }

}