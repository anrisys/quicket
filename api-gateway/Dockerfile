FROM nginx:alpine3.22-perl

# Remove default configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create log directory and set permissions
RUN mkdir -p /var/log/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R 755 /var/log/nginx 

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \ 
    CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# Expose port
EXPOSE 80

# Start nginx in foreground (tell docker to not stop nginx)
CMD [ "nginx", "-g", "daemon off" ]