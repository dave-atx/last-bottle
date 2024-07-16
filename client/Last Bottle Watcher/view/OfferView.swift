//
//  OfferView.swift
//  Last Bottle Watcher
//
//  Created by Marquard, Dave on 7/15/24.
//

import SwiftUI

struct OfferView: View {
    var offer: Offer
    
    var body: some View {
        VStack {
            AsyncImage(url: offer.imageURL) { image in
                image
                    .resizable()
                    .scaledToFit()
            } placeholder: {
                ProgressView()
            }
            .padding(.all)
            .frame(width: 200, height: 200)
                
            Text(offer.name)
                .font(.subheadline)
                .padding([.leading, .bottom, .trailing])
                
        }
    }
}

#Preview {
    OfferView(offer: Offer(id: "1234", name: "Great Napa Valley Chard 2024", price: 25, image: "https://s3.amazonaws.com/lastbottle/products/LBRDFJJ5-319332.jpg", purchaseUrl: "https://lastbottlewines.com"))
}
