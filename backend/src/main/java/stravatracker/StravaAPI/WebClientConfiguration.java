package stravatracker.StravaAPI;

import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpHeaders;
import org.springframework.web.reactive.function.client.ClientRequest;
import org.springframework.web.reactive.function.client.ExchangeFilterFunction;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.util.UriComponentsBuilder;
import reactor.core.publisher.Mono;

@Configuration
public class WebClientConfiguration {

    @Value("${strava.CLIENT_ID}")
    private String CLIENT_ID;

    @Value("${strava.CLIENT_SECRET}")
    private String CLIENT_SECRET;

    @Value("${strava.REFRESH_TOKEN}")
    private String REFRESH_TOKEN;


    @Bean
    public WebClient webClient() {
        return WebClient.builder()
                .baseUrl("https://www.strava.com/api/v3")
                .filter(oAuthFilter())
                .build();
    }


    private ExchangeFilterFunction oAuthFilter() {
        return ExchangeFilterFunction.ofRequestProcessor(clientRequest -> {
            ClientRequest authorizedRequest = ClientRequest.from(clientRequest)
                    .headers(headers -> headers.setBearerAuth(getAccessToken()))
                    .build();
            return Mono.just(authorizedRequest);
        });
    }

    private String getAccessToken() {

        HttpHeaders headers = new HttpHeaders();
        headers.set("Content-Type", "application/x-www-form-urlencoded");

        UriComponentsBuilder builder = UriComponentsBuilder
                .fromHttpUrl("https://www.strava.com/oauth/token")
                .queryParam("client_id", this.CLIENT_ID)
                .queryParam("client_secret", this.CLIENT_SECRET)
                .queryParam("refresh_token", this.REFRESH_TOKEN)
                .queryParam("grant_type", "refresh_token");

        WebClient.ResponseSpec responseSpec = WebClient
                .builder()
                .defaultHeaders(header -> header.addAll(headers))
                .build()
                .post()
                .uri(builder.toUriString())
                .retrieve();

        JsonNode responseBody = responseSpec.bodyToMono(JsonNode.class).block();

        if (responseBody != null) {
            System.out.println("Access token: " + responseBody.get("access_token").asText());
            return responseBody.get("access_token").asText();
        } else {
            throw new RuntimeException("Error obtaining access token");
        }


    }
}
