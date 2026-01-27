package api

// GraphQL query and mutation strings

const LoginMutation = `
mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
        authToken
        refreshToken
        __typename
    }
    version
    __typename
}
`

const RenewTokenMutation = `
mutation RenewToken($authToken: String!, $refreshToken: String!) {
    renewToken(authToken: $authToken, refreshToken: $refreshToken) {
        authToken
        refreshToken
    }
}
`

const MeQuery = `
query Me($siteReference: String) {
    me {
        id
        email
        countryCode
        advancedPaymentAmount(siteReference: $siteReference)
        treesCount
        hasInviteLink
        hasCO2Compensation
        createdAt
        updatedAt
        externalDetails {
            reference
            person {
                firstName
                lastName
            }
            contact {
                emailAddress
                phoneNumber
                mobileNumber
            }
            address {
                addressFormatted
                street
                houseNumber
                houseNumberAddition
                zipCode
                city
            }
        }
        smartCharging {
            isActivated
            provider
            isAvailableInCountry
        }
        smartTrading {
            isActivated
            isAvailableInCountry
        }
        websiteUrl
        customerSupportEmail
        reference
        connections(siteReference: $siteReference) {
            id
            connectionId
            EAN
            segment
            status
            contractStatus
            estimatedFeedIn
            firstMeterReadingDate
            lastMeterReadingDate
            meterType
            externalDetails {
                gridOperator
                address {
                    street
                    houseNumber
                    houseNumberAddition
                    zipCode
                    city
                }
                contract {
                    startDate
                    endDate
                    contractType
                    productName
                    tariffChartId
                }
            }
        }
    }
}
`

const UserSitesQuery = `
query UserSites {
    userSites {
        address {
            addressFormatted
        }
        addressHasMultipleSites
        deliveryEndDate
        deliveryStartDate
        firstMeterReadingDate
        lastMeterReadingDate
        propositionType
        reference
        segments
        status
    }
}
`

const MarketPricesQuery = `
query MarketPrices($date: String!, $resolution: PriceResolution!) {
    marketPrices(date: $date, resolution: $resolution) {
        averageElectricityPrices {
            averageMarketPrice
            averageMarketPricePlus
            averageAllInPrice
            perUnit
            isWeighted
            __typename
        }
        electricityPrices {
            from
            till
            resolution
            marketPrice
            marketPriceTax
            sourcingMarkupPrice
            energyTaxPrice
            marketPricePlus
            allInPrice
            perUnit
            __typename
        }
        gasPrices {
            from
            till
            resolution
            marketPrice
            marketPriceTax
            sourcingMarkupPrice
            energyTaxPrice
            marketPricePlus
            allInPrice
            perUnit
            __typename
        }
    }
}
`

const BelgiumMarketPricesQuery = `
query MarketPrices($date: String!) {
    marketPrices(date: $date) {
        electricityPrices {
            from
            till
            resolution
            marketPrice
            marketPriceTax
            sourcingMarkupPrice
            energyTaxPrice
            marketPricePlus
            allInPrice
            perUnit
            __typename
        }
        gasPrices {
            from
            till
            resolution
            marketPrice
            marketPriceTax
            sourcingMarkupPrice
            energyTaxPrice
            marketPricePlus
            allInPrice
            perUnit
            __typename
        }
        __typename
    }
}
`

const CustomerMarketPricesQuery = `
query MarketPrices($date: String!, $siteReference: String!) {
    customerMarketPrices(date: $date, siteReference: $siteReference) {
        id
        averageElectricityPrices {
            averageMarketPrice
            averageMarketPricePlus
            averageAllInPrice
            perUnit
            isWeighted
        }
        electricityPrices {
            id
            date
            from
            till
            resolution
            marketPrice
            marketPricePlus
            marketPriceTax
            sourcingMarkupPrice: consumptionSourcingMarkupPrice
            energyTaxPrice: energyTax
            allInPrice
            perUnit
            __typename
        }
        gasPrices {
            id
            date
            from
            till
            resolution
            marketPrice
            marketPricePlus
            marketPriceTax
            sourcingMarkupPrice: consumptionSourcingMarkupPrice
            energyTaxPrice: energyTax
            allInPrice
            perUnit
            __typename
        }
        __typename
    }
}
`

const MonthSummaryQuery = `
query MonthSummary($siteReference: String!) {
    monthSummary(siteReference: $siteReference) {
        _id
        actualCostsUntilLastMeterReadingDate
        expectedCostsUntilLastMeterReadingDate
        expectedCosts
        lastMeterReadingDate
        meterReadingDayCompleteness
        gasExcluded
        __typename
    }
    version
    __typename
}
`

const InvoicesQuery = `
query Invoices($siteReference: String!) {
    invoices(siteReference: $siteReference) {
        allInvoices {
            id
            invoiceDate
            startDate
            periodDescription
            totalAmount
            __typename
        }
        previousPeriodInvoice {
            id
            startDate
            periodDescription
            totalAmount
            __typename
        }
        currentPeriodInvoice {
            id
            startDate
            periodDescription
            totalAmount
            __typename
        }
        upcomingPeriodInvoice {
            id
            startDate
            periodDescription
            totalAmount
            __typename
        }
        __typename
    }
    __typename
}
`

const PeriodUsageAndCostsQuery = `
query PeriodUsageAndCosts($date: String!, $siteReference: String!) {
    periodUsageAndCosts(date: $date, siteReference: $siteReference) {
        _id
        gas {
            usageTotal
            costsTotal
            unit
            items {
                date
                from
                till
                usage
                costs
                unit
                __typename
            }
            __typename
        }
        electricity {
            usageTotal
            costsTotal
            unit
            items {
                date
                from
                till
                usage
                costs
                unit
                __typename
            }
            __typename
        }
        feedIn {
            usageTotal
            costsTotal
            unit
            items {
                date
                from
                till
                usage
                costs
                unit
                __typename
            }
            __typename
        }
        __typename
    }
    __typename
}
`

const EnodeChargersQuery = `
query EnodeChargers {
    enodeChargers {
        canSmartCharge
        chargeSettings {
            calculatedDeadline
            capacity
            deadline
            hourFriday
            hourMonday
            hourSaturday
            hourSunday
            hourThursday
            hourTuesday
            hourWednesday
            id
            initialCharge
            initialChargeTimestamp
            isSmartChargingEnabled
            isSolarChargingEnabled
            maxChargeLimit
            minChargeLimit
        }
        chargeState {
            batteryCapacity
            batteryLevel
            chargeLimit
            chargeRate
            chargeTimeRemaining
            isCharging
            isFullyCharged
            isPluggedIn
            lastUpdated
            powerDeliveryState
            range
        }
        id
        information {
            brand
            model
            year
        }
        interventions {
            description
            title
        }
        isReachable
        lastSeen
    }
}
`

const EnodeVehiclesQuery = `
query EnodeVehicles {
    enodeVehicles {
        canSmartCharge
        chargeSettings {
            calculatedDeadline
            deadline
            hourFriday
            hourMonday
            hourSaturday
            hourSunday
            hourThursday
            hourTuesday
            hourWednesday
            id
            isSmartChargingEnabled
            isSolarChargingEnabled
            maxChargeLimit
            minChargeLimit
        }
        chargeState {
            batteryCapacity
            batteryLevel
            chargeLimit
            chargeRate
            chargeTimeRemaining
            isCharging
            isFullyCharged
            isPluggedIn
            lastUpdated
            powerDeliveryState
            range
        }
        id
        information {
            brand
            model
            vin
            year
        }
        interventions {
            description
            title
        }
        isReachable
        lastSeen
    }
}
`

const SmartBatteriesQuery = `
query SmartBatteries {
    smartBatteries {
        brand
        capacity
        createdAt
        externalReference
        id
        maxChargePower
        maxDischargePower
        provider
        updatedAt
        __typename
    }
}
`

const SmartBatteryDetailsQuery = `
query SmartBattery($deviceId: String!) {
    smartBattery(deviceId: $deviceId) {
        brand
        capacity
        id
        settings {
            batteryMode
            imbalanceTradingStrategy
            selfConsumptionTradingAllowed
        }
    }
    smartBatterySummary(deviceId: $deviceId) {
        lastKnownStateOfCharge
        lastKnownStatus
        lastUpdate
        totalResult
    }
}
`

const SmartBatterySessionsQuery = `
query SmartBatterySessions($startDate: String!, $endDate: String!, $deviceId: String!) {
    smartBatterySessions(
        startDate: $startDate
        endDate: $endDate
        deviceId: $deviceId
    ) {
        deviceId
        fairUsePolicyVerified
        periodEndDate
        periodEpexResult
        periodFrankSlim
        periodImbalanceResult
        periodStartDate
        periodTotalResult
        periodTradeIndex
        periodTradingResult
        sessions {
            cumulativeResult
            date
            result
            status
            tradeIndex
        }
    }
}
`
