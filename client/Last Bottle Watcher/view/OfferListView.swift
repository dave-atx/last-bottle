//
//  OfferListView.swift
//  Last Bottle Watcher
//
//  Created by Marquard, Dave on 7/15/24.
//

import SwiftUI

struct OfferListView: View {
    @Environment(OfferModel.self) private var offerModel
    
    var body: some View {
        List(offerModel.offers) { offer in
            OfferView(offer: offer)
        }
        .task {
            try! await offerModel.refreshOffers()
        }
        .refreshable {
            try! await offerModel.refreshOffers()
        }
    }
}


#Preview {
    let offers = [
        Offer(id: "1234", name: "Great Napa Valley Chard 2024", price: 25, image: "https://s3.amazonaws.com/lastbottle/products/LBRDFJJ5-319332.jpg", purchaseUrl: "https://lastbottlewines.com"),
        Offer(id: "1234", name: "Great Napa Valley Chard 2024", price: 25, image: "https://s3.amazonaws.com/lastbottle/products/LBRDFJJ5-319332.jpg", purchaseUrl: "https://lastbottlewines.com")
    ]
    var offerModel = OfferModel()
    offerModel.offers = offers
    return OfferListView()
        .environment(offerModel)
}
