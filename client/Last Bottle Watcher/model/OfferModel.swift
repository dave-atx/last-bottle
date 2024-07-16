//
//  OfferModel.swift
//  Last Bottle Watcher
//
//  Created by Marquard, Dave on 7/15/24.
//

import Foundation

@Observable class OfferModel {
    let OFFERS_URL = URL(string: "https://lb.marquard.org/api/v1/offers")!
    var offers: [Offer] = []
    
    func refreshOffers() async throws {
        let (data, _) = try await URLSession.shared.data(from: OFFERS_URL)
        offers = try JSONDecoder().decode([Offer].self, from: data)
    }
}
