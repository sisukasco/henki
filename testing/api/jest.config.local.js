module.exports = {
    preset: 'ts-jest',
    testEnvironment: "node",
    testMatch:["**/tests/**/*.[jt]s"],
    testTimeout: 30000,
    "globals": {
        "test_config": {
            endpoint: "http://local.ratufa.io:3131",
            aud: "local.ratufa.io",
            emails: {
                endpoint: "http://local.ratufa.io:8025"
            }
        }
    }
}