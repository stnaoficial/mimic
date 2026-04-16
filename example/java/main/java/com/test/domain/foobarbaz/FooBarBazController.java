package com.test.domain.foobarbaz;

import org.springframework.web.bind.annotation.*;

@RestController
@Path("/api/v1/foo-bar-baz")
public interface FooBarBazController {

    @PostMapping()
    public String create(@RequestBody FooBarBazRequest fooBarBazRequest) {
        return "FooBarBaz created"
    }

    @GetMapping("/{id}")
    public String read(@PathVariable String fooBarBazId) {
        return "FooBarBaz found"
    }

}