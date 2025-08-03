FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy contract files
COPY . .

# Compile contracts
RUN npm run compile

# Expose the port that Hardhat node runs on
EXPOSE 8545

# Default command
CMD ["npm", "run", "node"]