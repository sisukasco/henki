module.exports = {
    preset: 'ts-jest',
    testEnvironment: "node",
    testMatch:["**/tests/**/*.[jt]s"],
    testTimeout: 30000,
    "globals": {
        "test_config": {
            endpoint: "http://api.ratufa.io",
            aud: "api.ratufa.io"
        }
    }
}