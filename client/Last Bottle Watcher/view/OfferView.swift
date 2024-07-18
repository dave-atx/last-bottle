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
        HStack {
            AsyncImage(url: offer.imageURL) { image in
                image
                    .resizable()
                    .scaledToFill()
                    .clipped()
                    .frame(width: 75, height: 150)
            } placeholder: {
                ProgressView()
            }
            .frame(width: 75, height: 150)
            .padding(.trailing)
            
            Text(offer.name)
                .font(.subheadline)
                .padding(.leading)
                .frame(maxHeight: 150)
            
            Spacer()
        }
        
    }
}

#Preview {
    OfferView(offer: Offer(id: "1234", name: "Great Napa Valley Chard 2024", price: 25, image: "https://s3.amazonaws.com/lastbottle/products/LBRDFJJ5-319332.jpg"))
}
