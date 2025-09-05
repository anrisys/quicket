# quicket/api-gateway/Dockerfile
FROM nginx:alpine-perl

# Remove default configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom configuration
COPY nginx.conf /etc/nginx/nginx.conf
# COPY includes/ /etc/nginx/includes/

# Expose port
EXPOSE 80

# Health check
# HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
#     CMD curl -f http://localhost/health || exit 1