package com.ronilsonalves.invoiceservice.config.feign;

import feign.RequestInterceptor;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.client.*;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.registration.ClientRegistration;
import org.springframework.security.oauth2.client.registration.ClientRegistrationRepository;

@Configuration
@RequiredArgsConstructor
public class FeignConfiguration {

    private static final String KEYCLOAK_REGISTRATION_ID = "keycloak-registration";

    private final OAuth2AuthorizedClientService clientService;
    private final ClientRegistrationRepository registrationRepository;

    @Bean
    public RequestInterceptor requestInterceptor() {
        ClientRegistration clientRegistration = registrationRepository.findByRegistrationId(KEYCLOAK_REGISTRATION_ID);
        OAAuth2ClientCredentialsFeignManager credentialsFeignMananger =
                new OAAuth2ClientCredentialsFeignManager(authorizedClientManager(),clientRegistration);
        return requestInterceptor -> {
            requestInterceptor.header("Authorization", "Bearer " + credentialsFeignMananger.getAccessToken());
        };
    }

    @Bean
    public OAuth2AuthorizedClientManager authorizedClientManager () {
        OAuth2AuthorizedClientProvider authorizedClientProvider = OAuth2AuthorizedClientProviderBuilder
                .builder()
                .clientCredentials()
                .build();

        AuthorizedClientServiceOAuth2AuthorizedClientManager authorizedClientManager =
                new AuthorizedClientServiceOAuth2AuthorizedClientManager(registrationRepository, clientService);

        authorizedClientManager.setAuthorizedClientProvider(authorizedClientProvider);

        return authorizedClientManager;
    }
}
