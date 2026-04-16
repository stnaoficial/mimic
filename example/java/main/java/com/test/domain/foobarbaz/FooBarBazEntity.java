package com.test.domain.foobarbaz;

import com.test.domain.FooBarBazBuilder;

public class FooBarBazEntity {

    private final String id;

    public FooBarBazEntity(FooBarBazBuilder fooBarBazBuilder) {
        this.id = builder.id;
    }

    public String getId() {
        return this.id;
    }

}