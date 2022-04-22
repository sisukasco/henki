module.exports = {
    preset: 'ts-jest',
    testEnvironment: "node",
    testMatch:["**/tests/**/*.[jt]s"],
    testTimeout: 30000,
    "globals": {
        "test_config": {
            endpoint: "http://localhost:3131",
            aud: "localhost",
            emails: {
                endpoint: "http://localhost:8025"
            }
        }
    }
}