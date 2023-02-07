package com.ronilsonalves.invoiceservice.config.oauth;


import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.boot.configurationprocessor.json.JSONException;
import org.springframework.boot.configurationprocessor.json.JSONObject;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.security.web.access.AccessDeniedHandler;

import java.io.IOException;
import java.time.LocalDateTime;

public class CustomAuthErrorResponse implements AuthenticationEntryPoint, AccessDeniedHandler {

    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response, AuthenticationException authException) throws IOException, ServletException {
        JSONObject jsonObject = new JSONObject();
        response.setContentType("application/json;charset=UTF-8");
        response.setStatus(401);
        buildResponse(response, jsonObject,"User must be authenticated to access this resource");
    }

    @Override
    public void handle(HttpServletRequest request, HttpServletResponse response, AccessDeniedException accessDeniedException) throws IOException, ServletException {
        JSONObject jsonObject = new JSONObject();
        response.setContentType("application/json;charset=UTF-8");
        response.setStatus(403);
        buildResponse(response, jsonObject, "User has not permission to access this resource");
    }

    private void buildResponse(HttpServletResponse response, JSONObject jsonObject, String message) throws IOException {
        try {
            jsonObject.put("timestamp", LocalDateTime.now());
            jsonObject.put("status",response.getStatus());
            jsonObject.put("message",message);
            response.getWriter().write(jsonObject.toString());
        } catch (JSONException e) {
            throw new RuntimeException(e);
        }
    }
}
