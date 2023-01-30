package com.ronilsonalves.invoiceservice.config;

import io.swagger.v3.oas.annotations.enums.SecuritySchemeType;
import io.swagger.v3.oas.annotations.security.OAuthFlow;
import io.swagger.v3.oas.annotations.security.OAuthFlows;
import io.swagger.v3.oas.models.parameters.HeaderParameter;
import io.swagger.v3.oas.models.security.SecurityScheme;
import io.swagger.v3.oas.models.Components;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import io.swagger.v3.oas.models.security.SecurityRequirement;
import org.springdoc.core.customizers.OpenApiCustomizer;
import org.springdoc.core.models.GroupedOpenApi;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
@io.swagger.v3.oas.annotations.security.SecurityScheme(name = "oauth_sec", type = SecuritySchemeType.OAUTH2, flows = @OAuthFlows(authorizationCode =
@OAuthFlow(
        authorizationUrl = "http://localhost:8090/realms/GoLangInSpringCloud/protocol/openid-connect/auth",
        refreshUrl = "http://localhost:8090/realms/GoLangInSpringCloud/protocol/openid-connect/token",
        tokenUrl = "http://localhost:8090/realms/GoLangInSpringCloud/protocol/openid-connect/token"
)))
public class SwaggerConfig {

    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
                .info(new Info().title("Invoice API")
                        .description("API for invoicing service")
                        .version("1")
                        .contact(new Contact().name("Ronilson Alves")
                                .email("hello@ronilsonalves")
                                .url("https://github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture"))
                        .license(new License().name("MIT License").url("https://github.com/ronilsonalves/GoLang-in-a" +
                                "-spring-cloud-architecture/blob/main/LICENSE")));
    }

    @Bean
    public GroupedOpenApi invoicesAPI() {
        return GroupedOpenApi.builder()
                .group("Invoices")
                .displayName("Invoices API")
                .pathsToMatch("/invoices/**")
                .build();
    }

    @Bean
    public OpenApiCustomizer openApiCustomizer() {
        return openApi -> openApi.getPaths().values().stream()
                .flatMap(pathItem -> pathItem.readOperations().stream())
                .forEach(operation -> operation.addParametersItem(new HeaderParameter()
                        .$ref("#/components/parameters/Version")));
    }
}
