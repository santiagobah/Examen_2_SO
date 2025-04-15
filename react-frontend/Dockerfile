# ---- Etapa de Construcción (Build Stage) ----
FROM node:18-alpine AS build

WORKDIR /app

# Copiar package.json y package-lock.json (o yarn.lock)
COPY package*.json ./

# Instalar dependencias
RUN npm install --verbose

# Copiar el resto de los archivos de la aplicación
COPY . .

# Construir la aplicación para producción
RUN npm run build

# ---- Etapa de Producción (Serve Stage) ----
FROM nginx:stable-alpine

# Copiar los archivos estáticos construidos desde la etapa 'build' al directorio web de Nginx
COPY --from=build /app/build /usr/share/nginx/html

# Copiar la configuración personalizada de Nginx (la crearemos a continuación)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Exponer el puerto 80 (puerto por defecto de Nginx)
EXPOSE 80

# Comando para iniciar Nginx en primer plano
CMD ["nginx", "-g", "daemon off;"]
